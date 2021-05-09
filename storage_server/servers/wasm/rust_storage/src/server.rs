// cargo run --bin storage-server
use tonic::{transport::Server, Request, Response, Status};
use wasi_cap_std_sync::WasiCtxBuilder;
use wasmtime;
use wasmtime_wasi::Wasi;

use anyhow::Result;
use prost::{bytes::BytesMut, Message};

use proto::storage_server::{Storage, StorageServer};
use proto::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};

use std::collections::HashMap;
use std::fs::File;
use std::path::Path;

use std::thread;
use tokio::sync::{mpsc, oneshot};

use cap_std::fs::Dir;

// IP is used to choose the IP of the server
const IP: &str = "152.94.162.17:50051";

pub mod proto {
    tonic::include_proto!("proto"); // The string specified here must match the proto package name
}

pub struct MyStorage {
    handle: WasmHandle,
}

impl MyStorage {
    pub fn new(handle: WasmHandle) -> Self {
        MyStorage { handle }
    }
}

#[tonic::async_trait]
impl Storage for MyStorage {
    // Accept request of type Message
    async fn read(&self, request: Request<ReadRequest>) -> Result<Response<ReadResponse>, Status> {
        let res = request.into_inner();
        let replay = self.handle.get_read_response(res).await;

        Ok(Response::new(replay))
    }

    async fn write(
        &self,
        request: Request<WriteRequest>,
    ) -> Result<Response<WriteResponse>, Status> {
        let res = request.into_inner();
        let replay = self.handle.get_write_response(res).await;

        Ok(Response::new(replay))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let path = Path::new("./data");
    let file = match File::open(&path) {
        Err(why) => panic!("coudn't open file: {}", why),
        Ok(file) => {
            println!("Could open .");
            file
        }
    };
    let dir: Dir = unsafe { Dir::from_std_file(file) };

    let handle: WasmHandle = WasmHandle::new(dir);

    // this is just a precoution to see if the WasmActor is ready
    let status = handle.ready_to_use().await;
    if status == 1 {
        println!("WasmActor ready")
    }

    // ---------

    let addr = IP.parse()?;
    let store_server = MyStorage::new(handle);
    println!("Server is running at {}\n", IP);

    Server::builder()
        .add_service(StorageServer::new(store_server))
        .serve(addr)
        .await?;

    Ok(())
}

// ActorRequest enum is used to specify how you want to interact with the actor
enum ActorRequest {
    // ReadyToUse is used to send a message when the actor has been initialised
    ReadyToUse {
        respond_to: oneshot::Sender<i32>,
    },
    Read {
        respond_to: oneshot::Sender<ReadResponse>,
        request: ReadRequest,
    },
    Write {
        respond_to: oneshot::Sender<WriteResponse>,
        request: WriteRequest,
    },
}

// WasmActor struct takes care of the Wasm instance and Wasm exported functions
struct WasmActor {
    receiver: mpsc::Receiver<ActorRequest>,
    funcs: HashMap<String, wasmtime::Func>,
    memory: wasmtime::Memory,
}

// WasmActor implementation
impl WasmActor {
    // new initializes the WasmActor, initializing the Wasm instance
    fn new(receiver: mpsc::Receiver<ActorRequest>, dir: Dir) -> Self {
        let engine = wasmtime::Engine::default();
        let store = wasmtime::Store::new(&engine);
        let mut linker = wasmtime::Linker::new(&store);

        // configurations
        let cx1 = WasiCtxBuilder::new();
        let cx1 = cx1.preopened_dir(dir, ".").expect("error opeing ");
        let cx1 = cx1.build().expect("Problem with wasctx");

        // link WASI
        let wasi = Wasi::new(&store, cx1);
        wasi.add_to_linker(&mut linker).expect("");

        // create the WebAssembly-module
        let module =
            wasmtime::Module::from_file(store.engine(), "../wasm_module/storage_application.wasm")
                .expect("problem creating the module");
        let instance = linker
            .instantiate(&module)
            .expect("problem creating an instance");

        // load and execute the _initialize function so that wasm gets access to the data folder
        let initialize = instance
            .get_func("_initialize")
            .expect("export wasn't a function");
        match initialize.call(&[]) {
            Ok(_result) => (),
            Err(trap) => {
                panic!("execution of initialize in a wasm trap: {}", trap);
            }
        };

        // export functions and memory from the WebAssemblt module
        let w_alloc = instance
            .get_func("new_alloc")
            .expect("export wasn't a function");
        let w_dealloc = instance
            .get_func("new_dealloc")
            .expect("export wasn't a function");
        let w_get_len = instance
            .get_func("get_response_len")
            .expect("export wasn't a function");
        // --------------------
        let w_write = instance
            .get_func("store_data")
            .expect("export wasn't a function");
        let w_read = instance
            .get_func("read_data")
            .expect("export wasn't a function");
        let mem = instance
            .get_memory("memory")
            .expect("memory export did not go well");

        // Store the funcs in the Actor struct
        let mut map: HashMap<String, wasmtime::Func> = HashMap::new();
        map.insert("alloc".to_string(), w_alloc);
        map.insert("dealloc".to_string(), w_dealloc);
        map.insert("get_len".to_string(), w_get_len);
        map.insert("write".to_string(), w_write);
        map.insert("read".to_string(), w_read);
        WasmActor {
            receiver,
            funcs: map,
            memory: mem,
        }
    }

    // handle_request handles any type of request. Request types behaviour has to be specified
    fn handle_request(&mut self, msg: ActorRequest) {
        match msg {
            ActorRequest::ReadyToUse { respond_to } => {
                // The `let _ =` ignores any errors when sending.
                //
                // This can happen if the `select!` macro is used
                // to cancel waiting for the response.
                let _ = respond_to.send(1);
            }
            ActorRequest::Read {
                respond_to,
                request,
            } => {
                let mut buf = BytesMut::with_capacity(500);
                request.encode(&mut buf).unwrap();
                let bytes_vec: Vec<u8> = buf.to_vec();
                let result = self.call_wasm("read", bytes_vec);

                let buf = &result[..];

                let response: proto::ReadResponse;
                response = prost::Message::decode(buf).unwrap();

                let _ = respond_to.send(response);
            }
            ActorRequest::Write {
                respond_to,
                request,
            } => {
                let mut buf = BytesMut::with_capacity(500);
                request.encode(&mut buf).unwrap();
                let bytes_vec: Vec<u8> = buf.to_vec();
                let result = self.call_wasm("write", bytes_vec);
                let buf = &result[..];

                let response: proto::WriteResponse = prost::Message::decode(buf).unwrap();

                let _ = respond_to.send(response);
            }
        }
    }
    // copy_to_memory handles the copy of serialized data to the Wasm's memory
    fn copy_to_memory(&mut self, data: Vec<u8>) -> (i32, i32) {
        let data = &data[..];
        let size = data.len();
        let alloc = self.funcs["alloc"]
            .get1::<i32, i32>()
            .expect("error converting alloc function");

        // allocate memory in wasm
        let ptr = alloc(size as i32).expect("something went wrong calling alloc");

        // write to memory
        let result = self.memory.write(ptr as usize, data);
        match result {
            Ok(result) => result,
            Err(e) => panic!("Error at write {}", e),
        };
        (ptr, size as i32)
    }

    // call_wasm handles the actuall wasm function call, and takes care of all calls to alloc/dialloc in the wasm instance
    fn call_wasm(&mut self, f_name: &str, data: Vec<u8>) -> Vec<u8> {
        let (ptr, len) = self.copy_to_memory(data);

        let func = self.funcs[f_name]
            .get2::<i32, i32, i32>()
            .expect("error converting `f` function");

        let w_deadlloc = self.funcs["dealloc"]
            .get2::<i32, i32, ()>()
            .expect("error converting dealloc function");

        let get_len = self.funcs["get_len"]
            .get0::<i32>()
            .expect("error converting get_len function");

        let res_ptr = func(ptr, len).expect("something went wrong calling the `f` function");

        // deallocate request protobuf message
        let _ = w_deadlloc(ptr, len);

        let result_len = get_len().expect("soemthing went wrong calling the get_len function");

        // create a buffer
        let mut buf: Vec<u8> = vec![0_u8; result_len as usize];
        //let mut buf = Vec::with_capacity(result_len as usize);
        let b: &mut [u8] = &mut buf[..];

        let write_result = self.memory.read(res_ptr as usize, b);

        match write_result {
            Ok(result) => result,
            Err(e) => panic!("Error at write {}", e),
        };

        // deallocate response protobuf message
        let _ = w_deadlloc(res_ptr, result_len);
        buf
    }
}

// run_my_actor runs the actor until it fails or program ends
fn run_my_actor(mut actor: WasmActor) {
    while let Some(msg) = actor.receiver.blocking_recv() {
        actor.handle_request(msg);
    }
}

#[derive(Clone)]
pub struct WasmHandle {
    sender: mpsc::Sender<ActorRequest>,
}

impl WasmHandle {
    pub fn new(dir: Dir) -> Self {
        let (sender, receiver) = mpsc::channel(8);
        thread::spawn(move || {
            let actor = WasmActor::new(receiver, dir);
            run_my_actor(actor);
        });
        Self { sender }
    }

    pub async fn ready_to_use(&self) -> i32 {
        let (send, recv) = oneshot::channel();
        let msg = ActorRequest::ReadyToUse { respond_to: send };

        // Ignore send errors. If this send fails, so does the
        // recv.await below. There's no reason to check for the
        // same failure twice.
        let _ = self.sender.send(msg).await;
        recv.await.expect("Actor task has been killed")
    }

    pub async fn get_write_response(&self, request: WriteRequest) -> WriteResponse {
        let (send, recv) = oneshot::channel();
        let msg = ActorRequest::Write {
            respond_to: send,
            request,
        };
        let _ = self.sender.send(msg).await;
        recv.await.expect("Actor task has been killed")
    }

    pub async fn get_read_response(&self, request: ReadRequest) -> ReadResponse {
        let (send, recv) = oneshot::channel();
        let msg = ActorRequest::Read {
            respond_to: send,
            request,
        };
        let _ = self.sender.send(msg).await;
        recv.await.expect("Actor task has been killed")
    }
}
