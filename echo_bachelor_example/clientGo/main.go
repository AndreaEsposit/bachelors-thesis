package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/AndreaEsposit/practice/echo_bachelor_example/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	check(err)

	client := pb.NewEchoClient(conn)

	fmt.Println("Exit/exit' to exit this program")

	for {
		reader := bufio.NewReader(os.Stdout)
		fmt.Print("Message to send: ")
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)

		if text == "exit" || text == "Exit" {
			break
		} else if text == "" {
			continue
		}

		message := &pb.Message{Content: text}

		returnMessage, err := client.Send(context.Background(), message)
		check(err)

		fmt.Printf("Recived this from server: '%v'\n", string(returnMessage.Content))
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
