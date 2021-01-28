extern crate protoc_rust;

//use protoc_rust::Customize;

fn main() {
    protoc_rust::Codegen::new()
        .out_dir("src/")
        .inputs(&["protos/echo.proto"])
        .include("protos")
        .run()
        .expect("protoc");
}
