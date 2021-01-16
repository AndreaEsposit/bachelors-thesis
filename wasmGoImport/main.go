package main

import (
	"fmt"
	"time"

	"github.com/bytecodealliance/wasmtime-go"
)

func main() {
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()

	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, _ := wasmtime.NewModuleFromFile(engine, "wasm_test.wasm")

	instance, _ := wasmtime.NewInstance(store, module, []*wasmtime.Extern{})

	age := instance.GetExport("age").Func()
	result, _ := age.Call(2021, 1998)

	fmt.Printf("You are: %v years old \n", result)
}
