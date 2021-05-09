# WebAssembly for systemprogramming
We have created a generalized WebAssembly+gRPC template that can quickly and efficiently be utilized to create an application for developers aiming to increase diversity thus, security in distributed applications. Following this template, a developer should be able to write highly reusable code that could be used in a variety of distributed applications.

Two applications can be found in this repository. While the echo application was developed to test Wasm's compatibility with gRPC, the storage application was created to be a practical example of our template. 



## Development
The project is developed using the following technologies:

* Golang 1.15+
* Python 3.16+
* Rust
* .NET 5
* gRPC
* WebAssembly
* Wasmtime
* Protocol Buffers

## Project structure
```
.
├── echo_server             # Echo server application files
├── storage_server          # Storage server application files
│   ├── benchmarks              # Benchmarking tools and utilities
│   ├── clients                 # Client applications
│   ├── proto                   # Compiled proto files
│   └── servers                 # Wasm and non-wasm server programs
├── wasm_modules            # WebAssembly/source files                 
├── go.mod
├── go.sum
└── README.md
```


## Authors
* Andrea Esposito
* John Marvin Cadacio
