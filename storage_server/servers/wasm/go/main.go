package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"github.com/bytecodealliance/wasmtime-go"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// IP is used to choose the IP of the server
const IP = "152.94.162.18:50051" // bbchain2=152.94.162.12

func main() {
	// ---------------------------------------------------------
	// initialize the WebAssembly module
	// stdout to print WASI text
	dir, err := ioutil.TempDir("", "out")
	check(err)
	defer os.RemoveAll(dir)
	stdoutPath := filepath.Join(dir, "stdout")

	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	linker := wasmtime.NewLinker(store)

	// configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdoutFile(stdoutPath)

	// pass access to this folder directory to the Wasm module
	err = wasiConfig.PreopenDir("./data", ".")
	check(err)

	// set the version to the same as in the WAT.
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	check(err)

	// link WASI
	err = linker.DefineWasi(wasi)
	check(err)

	// create the WebAssembly-module
	module, err := wasmtime.NewModuleFromFile(store.Engine, "../wasm_module/storage_application.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// execute the _initialize function to give wasm access to the data folder
	in := instance.GetExport("_initialize").Func()
	_, err = in.Call()
	if err != nil {
		panic(err)
	}

	// export functions and memory from the WebAssembly module
	funcs := make(map[string]*wasmtime.Func)
	funcs["alloc"] = instance.GetExport("new_alloc").Func()
	funcs["dealloc"] = instance.GetExport("new_dealloc").Func()
	funcs["get_len"] = instance.GetExport("get_response_len").Func()
	funcs["write"] = instance.GetExport("store_data").Func()
	funcs["read"] = instance.GetExport("read_data").Func()
	mem := instance.GetExport("memory").Memory()

	// -------------------------------------------------------------------------
	// initialize the grpc server
	server := NewStorageServer(funcs, mem, stdoutPath)
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
	stdout string
	memory *wasmtime.Memory
	funcs  map[string]*wasmtime.Func
	mu     sync.Mutex
	pb.UnimplementedStorageServer
	counter int
}

// NewStorageServer initializes an EchoServer
func NewStorageServer(funcs map[string]*wasmtime.Func, memory *wasmtime.Memory, stdout string) *StorageServer {
	return &StorageServer{
		funcs:  funcs,
		memory: memory,
		stdout: stdout,
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
	server.counter++
	fmt.Printf("Request #%v\n", server.counter)
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
	ptr, err := server.funcs["alloc"].Call(int32(len(data)))
	check(err)

	// casting pointer to int32
	ptr32 := ptr.(int32)

	// return raw memory backed by the WebAssembly memory as a byte slice
	buf := server.memory.UnsafeData()
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

	resPtr, err := server.funcs[fn].Call(ptr, len)
	check(err)
	resPtr32 := resPtr.(int32)

	// deallocate request protobuf message
	_, err = server.funcs["dealloc"].Call(ptr, len)
	check(err)

	resultLen, err := server.funcs["get_len"].Call()
	check(err)
	intResLen := resultLen.(int32)

	buf := server.memory.UnsafeData()
	response := make([]byte, int(intResLen))
	for i := range response {
		response[i] = buf[resPtr32+int32(i)]
	}

	// deallocate response protobuf message
	_, err = server.funcs["dealloc"].Call(resPtr32, intResLen)
	check(err)

	server.mu.Unlock()

	// unmarshalling
	if err := proto.Unmarshal(response, responseMessage); err != nil {
		log.Fatalln("Failed to parse message: ", err)
	}

	return responseMessage
}
