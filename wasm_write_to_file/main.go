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
	module, err := wasmtime.NewModuleFromFile(store.Engine, "new_rust_write.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// export functions and memory from the WebAssembly module

	// start := instance.GetExport("_start").Func()
	// _, err = start.Call()
	// if err != nil {
	// 	panic(err)
	// }

	in := instance.GetExport("_initialize").Func()
	_, err = in.Call()
	if err != nil {
		panic(err)
	}

	foo := instance.GetExport("foo").Func()
	_, err = foo.Call()
	check(err)

	// write := instance.GetExport("waWrite").Func()

	// res, err := write.Call()
	// check(err)
	// if res == true {
	// 	println("wrote to file")
	// } else {
	// 	println("didn't write to file")
	// }

	out, err := ioutil.ReadFile(stdoutPath)
	check(err)
	fmt.Print(string(out))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
