package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bytecodealliance/wasmtime-go"
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
	err = wasiConfig.PreopenDir(".", ".")
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

	alloc := instance.GetExport("new_alloc").Func()
	str := []byte("hello alloc function")
	ptr1, err := alloc.Call(len(str))
	//check(err)
	ptr132 := ptr1.(int32)

	mem := instance.GetExport("memory").Memory()
	buf := mem.UnsafeData()

	for i := range str {
		buf[ptr132+int32(i)] = str[i]
	}

	storeData := instance.GetExport("store_data").Func()
	_, err = storeData.Call(ptr132, len(str))
	check(err)

	dealloc := instance.GetExport("new_dealloc").Func()
	_, err = dealloc.Call(ptr132, len(str))
	check(err)

	rData := instance.GetExport("retrive_data").Func()
	newPtr, err := rData.Call()
	newPtr32 := newPtr.(int32)
	check(err)

	fmt.Printf("Pointer of 'retrive data' : %v\n\n", newPtr32)

	getlen := instance.GetExport("get_message_len").Func()
	nml, err := getlen.Call()
	check(err)
	newMessageLen := nml.(int32)

	newContent := make([]byte, newMessageLen)
	for i := range newContent {
		newContent[i] = buf[newPtr32+int32(i)]
	}

	fmt.Printf("THIS IS THE DATA from the file: '   %v   '\n", string(newContent))

	_, err = dealloc.Call(newPtr32, newMessageLen)
	check(err)

	newContent = make([]byte, newMessageLen)
	for i := range newContent {
		newContent[i] = buf[newPtr32+int32(i)]
	}

	fmt.Printf("THIS IS THE DATA AFTER DEALLOC: '   %v   '\n", string(newContent))

	rData2 := instance.GetExport("retrive_data2").Func()
	newPtr, err = rData2.Call()
	newPtr32 = newPtr.(int32)
	check(err)

	fmt.Printf("Pointer of 'retrive data2' : %v\n\n", newPtr32)

	nml, err = getlen.Call()
	check(err)
	newMessageLen = nml.(int32)

	newContent = make([]byte, newMessageLen)
	for i := range newContent {
		newContent[i] = buf[newPtr32+int32(i)]
	}

	fmt.Printf("THIS IS THE DATA AFTER DEALLOC: '   %v   '\n", string(newContent))

	out, err := ioutil.ReadFile(stdoutPath)
	check(err)
	fmt.Print(string(out))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
