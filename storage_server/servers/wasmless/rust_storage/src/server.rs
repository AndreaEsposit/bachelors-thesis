// cargo run --bin storage-server
use proto::storage_server::{Storage, StorageServer};
use proto::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};
use tonic::{transport::Server, Request, Response, Status};

use prost_types::Timestamp;
use serde_derive::{Deserialize, Serialize};
use std::path::Path;

pub mod proto {
    tonic::include_proto!("proto"); // The string specified here must match the proto package name
}

// Conent is used to read the file content
#[derive(Serialize, Deserialize, Debug)]
struct Content {
    nseconds: i32,
    seconds: i64,
    value: String,
}

pub struct MyStorage {
    sem: tokio::sync::Semaphore,
}

impl MyStorage {
    pub fn new(sem: tokio::sync::Semaphore) -> Self {
        MyStorage { sem }
    }
}

#[tonic::async_trait]
impl Storage for MyStorage {
    // Accept request of type Message
    async fn read(&self, request: Request<ReadRequest>) -> Result<Response<ReadResponse>, Status> {
        let request = request.into_inner();

        let mut file_path = "./data/".to_owned();
        file_path.push_str(&request.file_name);
        file_path.push_str(".json");

        let pathf = Path::new(&file_path);

        // get the semaphore permit
        let file_handle = self.sem.acquire().await;
        let data = tokio::fs::read(pathf).await;

        let response: ReadResponse;
        match data {
            Ok(data) => {
                let file_content: Content =
                    serde_json::from_slice(&data).expect("JSON was not well-formatted");

                let time: Option<Timestamp> = Some(Timestamp {
                    seconds: file_content.seconds,
                    nanos: file_content.nseconds,
                });

                // return response
                response = ReadResponse {
                    value: file_content.value,
                    ok: 1,
                    timestamp: time,
                };
            }
            Err(_e) => {
                // return response
                let time: Option<Timestamp> = Some(Timestamp {
                    seconds: 0,
                    nanos: 0,
                });
                // return response
                response = ReadResponse {
                    value: "".to_string(),
                    ok: 0,
                    timestamp: time,
                }
            }
        }
        drop(file_handle); // drop the lock
        Ok(Response::new(response))
    }

    async fn write(
        &self,
        request: Request<WriteRequest>,
    ) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();

        let mut file_path = "./data/".to_owned();
        file_path.push_str(&request.file_name);
        file_path.push_str(".json");

        let time = request.timestamp.expect("Error unwrapping the timestamp");

        let data = Content {
            seconds: time.seconds,
            nseconds: time.nanos,
            value: request.value,
        };

        let bdata = serde_json::to_vec_pretty(&data).unwrap();

        // acquire Semaphore permit
        let file_handle = self.sem.acquire().await;
        //let file_handle = self.mu.lock().expect("Mutex is poisoned");

        //write to file
        let e = tokio::fs::write(file_path, bdata).await;
        let e = match e {
            Ok(_result) => 1,
            Err(_e) => 0,
        };
        // drop the lock
        drop(file_handle);

        // return response
        let response: WriteResponse = WriteResponse { ok: e };

        Ok(Response::new(response))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "152.94.162.17:50051".parse()?;
    let mu = tokio::sync::Semaphore::new(1);
    let store_server = MyStorage::new(mu);
    println!("Server is running at {}\n", addr);

    Server::builder()
        .add_service(StorageServer::new(store_server))
        .serve(addr)
        .await?;

    Ok(())
}
