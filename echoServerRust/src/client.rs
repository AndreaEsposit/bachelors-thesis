// cargo run --bin echo-client
use proto::echo_client::EchoClient;
use proto::Message;
//use std::env; //used to get the env variables
use std::io::{self, Write};

pub mod proto {
    tonic::include_proto!("proto"); // Echo.proto package name = proto
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = EchoClient::connect("http://[::1]:50051").await?;

    //let args: Vec<String> = env::args().collect();

    print!("What are you thinking? ");
    io::stdout().flush().unwrap();
    let mut input = String::new();
    io::stdin().read_line(&mut input).unwrap();

    // read_line leaves a trailing newline, which trim will remove
    let command = input.trim();

    let request = tonic::Request::new(Message {
        content: command.as_bytes().into(),
    });

    let response = client.send(request).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
