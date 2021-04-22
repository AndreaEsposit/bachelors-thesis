package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"github.com/wasmerio/wasmer-go/wasmer"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// IP is used to choose the IP of the server
const IP = "localhost:50051" // bbchain2=152.94.162.12

func main() {
	// ---------------------------------------------------------
	// initialize the WebAssembly module

	// read wasm files
	wasmBytes, err := ioutil.ReadFile("../wasm_module/storage_application.wasm")
	check(err)

	var engine *wasmer.Engine
	// define engine and store
	if wasmer.IsCompilerAvailable(wasmer.LLVM) {
		config := wasmer.NewConfig()
		config.UseLLVMCompiler()
		engine = wasmer.NewEngineWithConfig(config)
		println("Using LLVM")
	} else {
		engine = wasmer.NewEngine()
	}

	store := wasmer.NewStore(engine)

	// Compiles the module
	module, err := wasmer.NewModule(store, wasmBytes)
	check(err)

	// configure WASI imports to write stdout into a file.
	wasiConfig, err := wasmer.NewWasiStateBuilder("storage").MapDirectory(".", "./data").Finalize()
	check(err)

	// set the version to the same as in the WAT.

	importObject, err := wasiConfig.GenerateImportObject(store, module)
	check(err)
	instance, err := wasmer.NewInstance(module, importObject)
	check(err)
	start, err := instance.Exports.GetFunction("_initialize")
	check(err)

	_, err = start()
	check(err)

	// export functions and memory from the WebAssembly module
	funcs := make(map[string]func(...interface{}) (interface{}, error))

	funcs["alloc"], err = instance.Exports.GetFunction("new_alloc")
	check(err)
	funcs["dealloc"], err = instance.Exports.GetFunction("new_dealloc")
	check(err)
	funcs["get_len"], err = instance.Exports.GetFunction("get_response_len")
	check(err)
	funcs["write"], err = instance.Exports.GetFunction("store_data")
	check(err)
	funcs["read"], err = instance.Exports.GetFunction("read_data")
	check(err)
	memory, err := instance.Exports.GetMemory("memory")
	check(err)

	// fmt.Println("Querying memory size...")
	// size := memory.Size()
	// fmt.Println("Memory size (pages):", size)
	// fmt.Println("Memory size (pages as bytes):", size.ToBytes())
	// fmt.Println("Memory size (bytes):", memory.DataSize())

	// -------------------------------------------------------------------------
	// initialize the grpc server
	server := NewStorageServer(funcs, memory)
	lis, err := net.Listen("tcp", server.port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterStorageServer(grpcServer, server)
	fmt.Printf("Server is running at %v.\n", server.port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

// StorageServer struct facilitates the managment of the server
type StorageServer struct {
	port   string
	memory *wasmer.Memory
	funcs  map[string]func(...interface{}) (interface{}, error)
	mu     sync.Mutex
	pb.UnimplementedStorageServer
}

// NewStorageServer initializes an EchoServer
func NewStorageServer(funcs map[string]func(...interface{}) (interface{}, error), memory *wasmer.Memory) *StorageServer {
	return &StorageServer{
		funcs:  funcs,
		memory: memory,
		port:   IP, //152.94.1.102:50051 (Pitter3)
	}
}

// Read will forward the protobuf message to the WebAssembly module and return what the module returns
func (server *StorageServer) Read(ctx context.Context, message *pb.ReadRequest) (*pb.ReadResponse, error) {
	wasmResponse := server.callWasm("read", message, &pb.ReadResponse{})
	return wasmResponse.(*pb.ReadResponse), nil
}

// Write will forward the protobuf message to the WebAssembly module and return what the module returns
func (server *StorageServer) Write(ctx context.Context, message *pb.WriteRequest) (*pb.WriteResponse, error) {
	wasmResponse := server.callWasm("write", message, &pb.WriteResponse{})
	return wasmResponse.(*pb.WriteResponse), nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// copyToMemory handles the copy of serialized data to the Wasm's memory
func (server *StorageServer) copyToMemory(data []byte) int32 {

	// allocate memory in wasm
	ptr, err := server.funcs["alloc"](int32(len(data)))
	check(err)

	// casting pointer to int32
	ptr32 := ptr.(int32)
	//fmt.Printf("This is the pointer %v\n", ptr32)

	// return raw memory backed by the WebAssembly memory as a byte slice
	buf := server.memory.Data()
	for i, v := range data {
		buf[ptr32+int32(i)] = v
	}
	// return the pointer
	return ptr32
}

// callWasm handles the actuall wasm function call, and takes care of all calls to alloc/dialloc in the wasm instance
func (server *StorageServer) callWasm(fn string, requestMessage proto.Message, responseMessage proto.Message) proto.Message {
	recivedBytes, err := proto.Marshal(requestMessage)
	check(err)

	// lock access to the server (extra security)
	server.mu.Lock()

	ptr := server.copyToMemory(recivedBytes)
	len := int32(len(recivedBytes))

	resPtr, err := server.funcs[fn](ptr, len)
	check(err)
	resPtr32 := resPtr.(int32)

	// deallocate request protobuf message
	_, err = server.funcs["dealloc"](ptr, len)
	check(err)

	resultLen, err := server.funcs["get_len"]()
	check(err)
	intResLen := resultLen.(int32)

	buf := server.memory.Data()
	// response := make([]byte, int(intResLen))
	// for i := range response {
	// 	response[i] = buf[resPtr32+int32(i)]
	// }

	// unmarshalling
	if err := proto.Unmarshal(buf[resPtr32:resPtr32+intResLen], responseMessage); err != nil {
		log.Fatalln("Failed to parse message: ", err)
	}

	// deallocate response protobuf message
	_, err = server.funcs["dealloc"](resPtr32, intResLen)
	check(err)

	server.mu.Unlock()

	return responseMessage
}
