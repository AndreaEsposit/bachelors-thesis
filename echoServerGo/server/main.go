package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/AndreaEsposit/practice/echoServerGo/proto"
	"google.golang.org/grpc"
)

func main() {
	server := NewEchoServer()
	lis, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterEchoServer(grpcServer, NewEchoServer())
	fmt.Printf("Server is running at :%v.\n", server.port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type EchoServer struct {
	port string
	pb.UnimplementedEchoServer
}

func NewEchoServer() *EchoServer {
	return &EchoServer{
		port: "8070",
	}
}

func (echo *EchoServer) Send(ctx context.Context, message *pb.Message) (*pb.Message, error) {

	fmt.Printf("Server recived: %v\n", message.Content)
	fmt.Println("Sending back")

	return message, nil

}
