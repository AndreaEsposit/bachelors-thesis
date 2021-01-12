package main

import (
	"context"
	"fmt"
	"log"

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
		fmt.Print("What are you thinking?")
		var answer string
		fmt.Scanln(&answer)

		if answer == "" {
			continue
		}

		message := &pb.Message{
			Content: []byte(answer),
		}

		returnMessage, err := client.Send(context.Background(), message)
		if err != nil {
			fmt.Println("Got an error: ", err)
		}

		fmt.Printf("Recived this from server: %v \n", string(returnMessage.Content))
	}
}
