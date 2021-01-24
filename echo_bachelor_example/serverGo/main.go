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

	pb "github.com/AndreaEsposit/practice/echo_bachelor_example/proto"
	"github.com/bytecodealliance/wasmtime-go"
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
	module, err := wasmtime.NewModuleFromFile(store.Engine, "../wasm/echo_server.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// export functions and memory from the WebAssembly module
	funcs := make(map[string]*wasmtime.Func)
	funcs["alloc"] = instance.GetExport("my_alloc").Func()
	funcs["dealloc"] = instance.GetExport("my_dealloc").Func()
	funcs["echo"] = instance.GetExport("echo").Func()
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
	fmt.Printf("Server is running at :%v.\n", server.port)

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
func (server *EchoServer) Send(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	server.mu.Lock()
	defer server.mu.Unlock()

	fmt.Printf("Server recived: '%v'\n", message.Content)

	ptr := server.copyMemory([]byte(message.Content))

	newPtr, err := server.funcs["echo"].Call(ptr, int32(len(message.Content)))
	check(err)
	newPtr32 := newPtr.(int32)

	// copy the bytes to a new buffer
	buf := server.memory.UnsafeData()
	newContent := make([]byte, len(message.Content))
	for i := range newContent {
		newContent[i] = buf[newPtr32+int32(i)]
	}

	// make new message
	returnMessage := pb.Message{Content: string(newContent)}

	// Deallocate memory in wasm
	_, err = server.funcs["dealloc"].Call(newPtr32, int32(len(message.Content)))

	// Print WASM stdout
	// out, err := ioutil.ReadFile(echo.stdout)
	// check(err)
	// fmt.Print(string(out))

	return &returnMessage, nil

}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (server *EchoServer) copyMemory(data []byte) int32 {

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
