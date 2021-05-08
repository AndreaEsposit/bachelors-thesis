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

	pb "github.com/AndreaEsposit/bachelors-thesis/echo_server/proto"
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

func main() {
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

	// set the version to the same as in the WAT.
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	check(err)

	// link WASI
	err = linker.DefineWasi(wasi)
	check(err)

	// create the WebAssembly-module
	module, err := wasmtime.NewModuleFromFile(store.Engine, "../wasm_module/echo_server.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// export functions and memory from the WebAssembly module
	funcs := make(map[string]*wasmtime.Func)
	funcs["alloc"] = instance.GetExport("new_alloc").Func()
	funcs["dealloc"] = instance.GetExport("new_dealloc").Func()
	funcs["echo"] = instance.GetExport("echo").Func()
	funcs["get_len"] = instance.GetExport("get_message_len").Func()
	mem := instance.GetExport("memory").Memory()
	// -------------------------------------------------------------------------
	// initialize the grpc server
	server := NewEchoServer(funcs, mem, stdoutPath)
	lis, err := net.Listen("tcp", server.port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterEchoServer(grpcServer, server)
	fmt.Printf("Server is running at %v.\n", server.port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

// EchoServer struct facilitates the managment of the server
type EchoServer struct {
	port   string
	stdout string
	memory *wasmtime.Memory
	funcs  map[string]*wasmtime.Func
	mu     sync.Mutex
	pb.UnimplementedEchoServer
}

// NewEchoServer initializes an EchoServer
func NewEchoServer(funcs map[string]*wasmtime.Func, memory *wasmtime.Memory, stdout string) *EchoServer {
	return &EchoServer{
		funcs:  funcs,
		memory: memory,
		stdout: stdout,
		port:   "localhost:50051",
	}
}

// Send is the function called by the clients
func (server *EchoServer) Send(ctx context.Context, message *pb.EchoMessage) (*pb.EchoMessage, error) {
	//fmt.Printf("Server recived: '%v'\n", message.Content)
	recivedBytes, err := proto.Marshal(message)
	check(err)

	server.mu.Lock()
	defer server.mu.Unlock()

	ptr := server.copyToMemory(recivedBytes)

	newPtr, err := server.funcs["echo"].Call(ptr, int32(len(recivedBytes)))
	check(err)
	newPtr32 := newPtr.(int32)

	nml, err := server.funcs["get_len"].Call()
	check(err)
	newMessageLen := nml.(int32)

	buf := server.memory.UnsafeData()
	// make new message
	returnMessage := &pb.EchoMessage{}
	if err := proto.Unmarshal(buf[newPtr32:newPtr32+newMessageLen], returnMessage); err != nil {
		log.Fatalln("Failed to parse message: ", err)
	}

	// Deallocate memory in wasm
	_, err = server.funcs["dealloc"].Call(ptr, int32(len(recivedBytes)))
	_, err = server.funcs["dealloc"].Call(newPtr32, newMessageLen)

	// Print WASM stdout
	// out, err := ioutil.ReadFile(echo.stdout)
	// check(err)
	// fmt.Print(string(out))
	return returnMessage, nil
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (server *EchoServer) copyToMemory(data []byte) int32 {

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
