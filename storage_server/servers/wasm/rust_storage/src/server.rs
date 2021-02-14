// cargo run --bin storage-server
use tonic::{transport::Server, Request, Response, Status};
use wasmtime;
use wasmtime_wasi::{Wasi, WasiCtxBuilder};

use anyhow::Result;
use prost_types;

use proto::storage_server::{Storage, StorageServer};
use proto::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};

use std::collections::HashMap;
use std::fs::File;
use std::path::Path;

use std::thread;
use tokio::sync::{mpsc, oneshot};

pub mod proto {
    tonic::include_proto!("proto"); // The string specified here must match the proto package name
}

#[derive(Debug, Default)]
pub struct MyStorage {}

#[tonic::async_trait]
impl Storage for MyStorage {
    // Accept request of type Message
    async fn read(&self, _request: Request<ReadRequest>) -> Result<Response<ReadResponse>, Status> {
        let replay = proto::ReadResponse {
            ok: 0,
            timestamp: Some(prost_types::Timestamp {
                seconds: 1,
                nanos: 42,
            }),
            value: "Hello".to_string(),
        };
        Ok(Response::new(replay))
    }

    async fn write(
        &self,
        _request: Request<WriteRequest>,
    ) -> Result<Response<WriteResponse>, Status> {
        let replay = proto::WriteResponse { ok: 0 };
        Ok(Response::new(replay))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let path = Path::new(".");
    let file = match File::open(&path) {
        Err(why) => panic!("coudn't open file: {}", why),
        Ok(file) => {
            println!("Could open .");
            file
        }
    };

    let handler: MyActorHandle = MyActorHandle::new(file);
    let id = handler.get_unique_id().await;

    println!("This is the id: {:?}", id);

    // ---------

    let addr = "127.0.0.1:50051".parse()?; // ? used insted of Match here (switch in rust)
    let store_server = MyStorage::default();
    println!("Server is running at 127.0.0.1:50051\n");

    Server::builder()
        .add_service(StorageServer::new(store_server))
        .serve(addr)
        .await?;

    Ok(())
}

struct MyActor {
    receiver: mpsc::Receiver<ActorMessage>,
    funcs: HashMap<String, wasmtime::Func>,
    n: i32,
}
enum ActorMessage {
    GetFuncs { respond_to: oneshot::Sender<i32> },
}

impl MyActor {
    fn new(receiver: mpsc::Receiver<ActorMessage>, file: File) -> Self {
        let engine = wasmtime::Engine::default();
        let store = wasmtime::Store::new(&engine);
        let mut linker = wasmtime::Linker::new(&store);

        let mut cx1 = WasiCtxBuilder::new();
        cx1.preopened_dir(file, ".");
        let cx1 = cx1.build().expect("Problem with wasctx");

        let wasi = Wasi::new(&store, cx1);
        wasi.add_to_linker(&mut linker).expect("");

        // wasiCtx.preopened_dir(dir: fs::File::, guest_path: P)

        let module = wasmtime::Module::from_file(store.engine(), "../wasm_module/write.wasm")
            .expect("problem creating the module");
        let instance = linker
            .instantiate(&module)
            .expect("problem creating an instance");

        let initialize = instance
            .get_func("_initialize")
            .expect("export wasn't a function");
        match initialize.call(&[]) {
            Ok(_result) => (),
            Err(trap) => {
                panic!("execution of initialize in a wasm trap: {}", trap);
            }
        };

        let write = instance
            .get_func("store_data")
            .expect("export wasn't a function");

        let mut map: HashMap<String, wasmtime::Func> = HashMap::new();
        map.insert("initialize".to_string(), initialize);
        map.insert("write".to_string(), write);
        MyActor {
            receiver,
            funcs: map,
            n: 2,
        }
    }

    fn handle_message(&mut self, msg: ActorMessage) {
        match self.funcs["write"].call(&[]) {
            Ok(_result) => (),
            Err(trap) => {
                panic!("execution of initialize in a wasm trap: {}", trap);
            }
        };
        match msg {
            ActorMessage::GetFuncs { respond_to } => {
                // The `let _ =` ignores any errors when sending.
                //
                // This can happen if the `select!` macro is used
                // to cancel waiting for the response.
                let _ = respond_to.send(self.n);
            }
        }
    }
}

fn run_my_actor(mut actor: MyActor) {
    while let Some(msg) = actor.receiver.blocking_recv() {
        actor.handle_message(msg);
    }
}

#[derive(Clone)]
pub struct MyActorHandle {
    sender: mpsc::Sender<ActorMessage>,
}

impl MyActorHandle {
    pub fn new(file: File) -> Self {
        let (sender, receiver) = mpsc::channel(8);
        thread::spawn(move || {
            let actor = MyActor::new(receiver, file);
            run_my_actor(actor);
        });
        Self { sender }
    }

    pub async fn get_unique_id(&self) -> i32 {
        let (send, recv) = oneshot::channel();
        let msg = ActorMessage::GetFuncs { respond_to: send };

        // Ignore send errors. If this send fails, so does the
        // recv.await below. There's no reason to check for the
        // same failure twice.
        let _ = self.sender.send(msg).await;
        recv.await.expect("Actor task has been killed")
    }
}
