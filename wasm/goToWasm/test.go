//tinygo build -opt=s -o program.wasm -wasm-abi=generic -target=wasi program.go
package main

func main() {

}

//export add
func add(x int, y int) int {
	total := 0
	total = x + y
	return total
}
