// cargo run --bin echo-server
use tonic::{transport::Server, Request, Response, Status};

use proto::echo_server::{Echo, EchoServer}; // Services
use proto::EchoMessage; // Messages

pub mod proto {
    tonic::include_proto!("proto"); // The string specified here must match the proto package name
}

#[derive(Debug, Default)]
pub struct MyEcho {}

#[tonic::async_trait]
impl Echo for MyEcho {
    // Accept request of type Message
    async fn send(&self, request: Request<EchoMessage>) -> Result<Response<EchoMessage>, Status> {
        // Return an instance of type Message

        let message: String = request.into_inner().content;
        println!("Server recived '{}'", message);

        let reply = proto::EchoMessage {
            content: message.into(), // We must use .into_inner() as the fields of gRPC requests and responses are private
        };

        Ok(Response::new(reply)) // Send back our formatted greeting
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?; // ? used insted of Match here (switch in rust)
    let echo_server = MyEcho::default();
    println!("Server is running at :localhost:50051\n");

    Server::builder()
        .add_service(EchoServer::new(echo_server))
        .serve(addr)
        .await?;

    Ok(())
}
