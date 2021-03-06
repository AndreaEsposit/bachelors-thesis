// cargo run --bin storage-server
use proto::storage_server::{Storage, StorageServer};
use proto::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};
use tonic::{transport::Server, Request, Response, Status};

use prost_types::Timestamp;
use serde_derive::{Deserialize, Serialize};
use serde_json::json;
use std::{fs::File, io, path::Path};

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

#[derive(Debug, Default)]
pub struct MyStorage {}

#[tonic::async_trait]
impl Storage for MyStorage {
    // Accept request of type Message
    async fn read(&self, request: Request<ReadRequest>) -> Result<Response<ReadResponse>, Status> {
        let request = request.into_inner();

        let mut file_path = "./data".to_owned();
        file_path.push_str(&request.file_name);
        file_path.push_str(".json");
        let pathf = Path::new(&file_path);
        let file = File::open(pathf);

        let response: ReadResponse;
        match file {
            Ok(file) => {
                let reader = io::BufReader::new(file);

                let file_content: Content =
                    serde_json::from_reader(reader).expect("JSON was not well-formatted");

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

        let data = json!({
        "seconds" : time.seconds,
        "nseconds": time.nanos,
        "value": request.value,});

        // write to file
        let write_result = write_to_file(file_path, data);
        let write_result = match write_result {
            Ok(_result) => 1,
            Err(_e) => 0,
        };

        // return response
        let response: WriteResponse = WriteResponse { ok: write_result };

        Ok(Response::new(response))
    }
}

fn write_to_file(file_path: String, data: serde_json::Value) -> Result<(), io::Error> {
    let file = File::create(file_path)?;
    let e = serde_json::to_writer_pretty(file, &data)?;
    Ok(e)
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "127.0.0.1:50051".parse()?;
    let store_server = MyStorage::default();
    println!("Server is running at 127.0.0.1:50051\n");

    Server::builder()
        .add_service(StorageServer::new(store_server))
        .serve(addr)
        .await?;

    Ok(())
}