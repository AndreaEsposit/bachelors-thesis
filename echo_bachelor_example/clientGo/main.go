package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/AndreaEsposit/practice/echoServerGo/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("[::1]:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewEchoClient(conn)

	for {
		reader := bufio.NewReader(os.Stdout)
		fmt.Println("Exit/exit' to exit this program")
		fmt.Print("What are you thinking? ")
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)

		if text == "" {
			continue
		} else if text == "exit" || text == "Exit" {
			break
		}

		message := &pb.Message{
			Content: []byte(text),
		}

		returnMessage, err := client.Send(context.Background(), message)
		if err != nil {
			fmt.Println("Got an error: ", err)
		}

		fmt.Printf("Recived this from server: %v \n", string(returnMessage.Content))
	}
}
