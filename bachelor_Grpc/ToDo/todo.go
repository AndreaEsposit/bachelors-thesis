package main

func main() {

}

// //export addTodo
// func addTodo(todo string, list []string) {
// 	fmt.Printf("You have %v ToDos\n", len(list))
// 	list = append(list, todo)
// }

//export deleteTodo
func deleteTodo(i int, list []string) []string {
	index := i - 1
	if index == 0 {
		list = list[1:]
	} else if index == len(list)-1 {
		list = list[:len(list)-1]
	} else {
		copy(list[index:], list[index+1:])
		return list[:len(list)-1]
	}
	return list
}

// //export getTodos
// func getTodos(list []string) {
// 	for i, todo := range list {
// 		fmt.Printf("- %d: %s\n", i, todo)
// 	}
// }

//export editTodo
func editTodo(new string, index int, list []string) []string {
	list[index] = new
	return list
}
