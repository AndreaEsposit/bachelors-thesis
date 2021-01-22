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
	// Stdout to print WASI text
	dir, err := ioutil.TempDir("", "out")
	check(err)
	defer os.RemoveAll(dir)
	stdoutPath := filepath.Join(dir, "stdout")

	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	linker := wasmtime.NewLinker(store)

	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdoutFile(stdoutPath)

	// Set the version to the same as in the WAT.
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	check(err)

	// Link WASI
	err = linker.DefineWasi(wasi)
	check(err)

	// Create our module
	module, err := wasmtime.NewModuleFromFile(store.Engine, "../wasm/echo_server.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// --------------------------------------------------------------

	// Initialize the grpc server
	server := NewEchoServer(instance, stdoutPath)
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

type EchoServer struct {
	port     string
	stdout   string
	instance *wasmtime.Instance
	mu       sync.Mutex
	pb.UnimplementedEchoServer
}

func NewEchoServer(instance *wasmtime.Instance, stdout string) *EchoServer {
	return &EchoServer{
		port:     "localhost:50051", //[::1]:50051
		instance: instance,
		stdout:   stdout,
	}
}

func (echo *EchoServer) Send(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	echo.mu.Lock()
	defer echo.mu.Unlock()
	mem := echo.instance.GetExport("memory").Memory()
	alloc := echo.instance.GetExport("my_alloc").Func()
	deAlloc := echo.instance.GetExport("my_dealloc").Func()
	fn := echo.instance.GetExport("echo").Func()

	fmt.Printf("Server recived: '%v'\n", string(message.Content))

	ptr := copyMemory([]byte(message.Content), alloc, mem)

	newPtr, err := fn.Call(ptr, int32(len(message.Content)))
	check(err)
	pointer := newPtr.(int32)

	// Copy the bytes to a new buffer
	buf := mem.UnsafeData()
	newContent := make([]byte, len(message.Content))
	for i := range newContent {
		newContent[i] = buf[pointer+int32(i)]
	}

	returnMessage := pb.Message{
		Content: string(newContent),
	}

	// Deallocate memory in wasm
	_, err = deAlloc.Call(pointer, int32(len(message.Content)))

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

func copyMemory(data []byte, alloc *wasmtime.Func, mem *wasmtime.Memory) int32 {
	// find size of data
	size := int32(len(data))

	// allocate memory in wasm
	ptr, err := alloc.Call(size)
	check(err)

	pointer := ptr.(int32)

	buf := mem.UnsafeData()
	for i, v := range data {
		buf[pointer+int32(i)] = v
	}
	// return the pointer
	return pointer
}
