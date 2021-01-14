package main

import (
	"fmt"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func main() {

	bytes, _ := wasm.ReadBytes("wasm_test.wasm")

	// Instantiates the WebAssembly module.
	instance, _ := wasm.NewInstance(bytes)
	defer instance.Close()

	// Gets the `age` exported function from the WebAssembly instance.
	age := instance.Exports["age"]

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	result, _ := age(2021, 1998)

	fmt.Printf("You are: %v years old \n", result)
}
