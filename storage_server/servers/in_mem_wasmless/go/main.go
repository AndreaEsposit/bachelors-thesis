package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// IP is used to choose the IP of the server
const IP = "152.94.162.12:50051" //152.94.162.31:50051 //bbchain21

// StorageServer is the representation of the server
type StorageServer struct {
	port string
	pb.UnimplementedStorageServer
	mu sync.Mutex
}

// Data rapresents the storage
type Data struct {
	Seconds  int64  `json:"seconds"`
	Nseconds int32  `json:"nseconds"`
	Value    string `json:"value"`
}

// STORAGE is the in memory variable version of the storage file
var STORAGE = Data{Seconds: 0, Nseconds: 0, Value: ""}

// NewStorageServer initializes an EchoServer
func NewStorageServer() *StorageServer {
	return &StorageServer{
		port: IP, // bbchain21
	}
}

func (server *StorageServer) Read(ctx context.Context, request *pb.ReadRequest) (*pb.ReadResponse, error) {
	filename := request.FileName

	var response = pb.ReadResponse{}

	if filename != "test" {
		timestamp := timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		}

		// return response
		response.Value = ""
		response.Timestamp = &timestamp
		response.Ok = 0

	} else {
		server.mu.Lock()
		timestamp := timestamppb.Timestamp{
			Seconds: STORAGE.Seconds,
			Nanos:   STORAGE.Nseconds,
		}

		// return response
		response.Value = STORAGE.Value
		response.Timestamp = &timestamp
		response.Ok = 1
		server.mu.Unlock()
	}
	return &response, nil
}

func (server *StorageServer) Write(ctx context.Context, request *pb.WriteRequest) (*pb.WriteResponse, error) {

	timestamp := request.Timestamp

	// in memory storage
	server.mu.Lock()
	STORAGE.Nseconds = timestamp.Nanos
	STORAGE.Seconds = timestamp.Seconds
	STORAGE.Value = request.Value
	server.mu.Unlock()

	// return response
	response := &pb.WriteResponse{
		Ok: 1,
	}

	return response, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// initialize the grpc server
func main() {
	server := NewStorageServer()
	lis, err := net.Listen("tcp", server.port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterStorageServer(grpcServer, server)
	fmt.Printf("Server is running at %v.\n", server.port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
