package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "github.com/AndreaEsposit/practice/wasm_write_to_file/proto"
	"github.com/bytecodealliance/wasmtime-go"
	"google.golang.org/protobuf/proto"
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
	err = wasiConfig.PreopenDir("./data", ".")
	check(err)

	// set the version to the same as in the WAT.
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	check(err)

	// link WASI
	err = linker.DefineWasi(wasi)
	check(err)

	// create the WebAssembly-module
	module, err := wasmtime.NewModuleFromFile(store.Engine, "storage_application.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// export functions and memory from the WebAssembly module

	in := instance.GetExport("_initialize").Func()
	_, err = in.Call()
	if err != nil {
		panic(err)
	}

	mem := instance.GetExport("memory").Memory()
	alloc := instance.GetExport("new_alloc").Func()
	dealloc := instance.GetExport("new_dealloc").Func()
	//write := instance.GetExport("store_data").Func()
	getResultLen := instance.GetExport("get_response_len").Func()
	readData := instance.GetExport("read_data").Func()

	// timep, err := ptypes.TimestampProto(time.Now())
	// writeM := pb.WriteRequest{FileName: "Wasm", Value: "Important string", Timestamp: timep}

	// dataBytes, err := proto.Marshal(&writeM)
	// check(err)

	// response := callFunction(write, getResultLen, alloc, dealloc, mem, dataBytes)

	// returnMessage := &pb.WriteResponse{}
	// if err := proto.Unmarshal(response, returnMessage); err != nil {
	// 	log.Fatalln("Failed to parse message: ", err)
	// }

	// r := returnMessage.GetOk()
	// if r == 1 {
	// 	fmt.Println("We managed")

	// } else {
	// 	fmt.Println("We fucked up")
	// }

	// timep, err = ptypes.TimestampProto(time.Now())
	// writeM = pb.WriteRequest{FileName: "Wasm", Value: "Important string", Timestamp: timep}

	// dataBytes, err = proto.Marshal(&writeM)
	// check(err)

	// response = callFunction(write, getResultLen, alloc, dealloc, mem, dataBytes)

	// returnMessage = &pb.WriteResponse{}
	// if err := proto.Unmarshal(response, returnMessage); err != nil {
	// 	log.Fatalln("Failed to parse message: ", err)
	// }

	// r = returnMessage.GetOk()
	// if r == 1 {
	// 	fmt.Println("We managed")

	// } else {
	// 	fmt.Println("We fucked up")
	// }

	readM := pb.ReadRequest{FileName: "asm"}
	dataBytes, err := proto.Marshal(&readM)
	check(err)
	response := callFunction(readData, getResultLen, alloc, dealloc, mem, dataBytes)

	returnM := &pb.ReadResponse{}
	if err := proto.Unmarshal(response, returnM); err != nil {
		log.Fatalln("Failed to parse message: ", err)
	}

	fmt.Printf("This is the answer: %v\n", returnM)

	out, err := ioutil.ReadFile(stdoutPath)
	check(err)
	fmt.Print(string(out))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func copyMemory(alloc *wasmtime.Func, memory *wasmtime.Memory, data []byte) int32 {

	// allocate memory in wasm
	ptr, err := alloc.Call(int32(len(data)))
	check(err)

	// casting pointer to int32
	ptr32 := ptr.(int32)

	// return raw memory backed by the WebAssembly memory as a byte slice
	buf := memory.UnsafeData()
	for i, v := range data {
		buf[ptr32+int32(i)] = v
	}
	// return the pointer
	return ptr32
}

func callFunction(fn, getSize, alloc, dealloc *wasmtime.Func, memory *wasmtime.Memory, data []byte) (response []byte) {
	ptr := copyMemory(alloc, memory, data)
	len := int32(len(data))

	resPtr, err := fn.Call(ptr, len)
	check(err)
	resPtr32 := resPtr.(int32)

	// deallocate request protobuf message
	_, err = dealloc.Call(ptr, len)
	check(err)

	resultLen, err := getSize.Call()
	check(err)
	intResLen := resultLen.(int32)

	buf := memory.UnsafeData()
	response = make([]byte, int(intResLen))
	for i := range response {
		response[i] = buf[resPtr32+int32(i)]
	}

	// deallocate response protobuf message
	_, err = dealloc.Call(resPtr32, intResLen)
	check(err)

	return response
}
