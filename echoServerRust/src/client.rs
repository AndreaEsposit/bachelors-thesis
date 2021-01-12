use proto::echo_client::EchoClient;
use proto::Message;
use std::env;

pub mod proto {
    tonic::include_proto!("proto"); // Echo.proto package name = proto
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = EchoClient::connect("http://[::1]:50051").await?;

    let args: Vec<String> = env::args().collect();

    let request = tonic::Request::new(Message {
        content: args[1].as_bytes().into(),
    });

    let response = client.send(request).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
