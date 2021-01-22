package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bytecodealliance/wasmtime-go"
)

func main() {
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
	module, err := wasmtime.NewModuleFromFile(store.Engine, "to_do_rust/target/wasm32-wasi/debug/to_do_rust.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	fn := instance.GetExport("greet").Func()
	addVec := instance.GetExport("array_sum").Func()
	alloc := instance.GetExport("my_alloc").Func()
	memory := instance.GetExport("memory").Memory()

	//fmt.Printf("Memory size: %v\n", memory.Size())
	//fmt.Printf("Memory datasize: %v\n", memory.DataSize())

	size1 := int32(len([]byte("Andrea")))
	size2 := int32(len([]byte{1, 2, 5}))

	strs := []string{"foo", "bar"}
	buffa := &bytes.Buffer{}
	gob.NewEncoder(buffa).Encode(strs)
	bs := buffa.Bytes()

	size3 := int32(len(bs))

	// //Allocate memomory
	ptr1, err := alloc.Call(size1)
	check(err)
	pointer, _ := ptr1.(int32)

	pt2, err := alloc.Call(size2)
	pointe2, _ := pt2.(int32)

	ptr3, err := alloc.Call(size3)
	pointer3, _ := ptr3.(int32)

	//fmt.Printf("New size: %v\n", memory.Size())
	//fmt.Printf("New datasize: %v\n", memory.DataSize())

	buf := memory.UnsafeData()
	for i, v := range []byte("Andrea") {
		buf[pointer+int32(i)] = v
	}

	// Use string func
	_, err = fn.Call(pointer, size1)
	check(err)

	// Allocate array on memory
	for i, v := range []byte{1, 2, 5} {
		buf[pointe2+int32(i)] = v
	}

	// Call array_sum
	sum, err := addVec.Call(pointe2, size2)
	fmt.Printf("The sum is: %d\n", sum)

	// Allocate array of string on memory
	for i, v := range bs {
		buf[pointer3+int32(i)] = v
	}
	_, err = addVec.Call(pointe2, size2)

	// Print WASM stdout
	out, err := ioutil.ReadFile(stdoutPath)
	check(err)
	fmt.Print(string(out))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
