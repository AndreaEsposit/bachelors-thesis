package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"google.golang.org/grpc"
)

type StorageServer struct {
	port string
	pb.UnimplementedStorageServer
}

type Data struct {
	seconds  int64
	nseconds int32
	value    string
}

// NewStorageServer initializes an EchoServer
func NewStorageServer() *StorageServer {
	return &StorageServer{
		port: "localhost:50051", //152.94.1.102:50051 (Pitter3)
	}
}

func (server *StorageServer) Read(ctx context.Context, request *pb.ReadRequest) (*pb.ReadResponse, error) {
	// filename := request.FileName
	// fpath := "./data/" + filename + ".json"
	// file, err := os.Open(fpath)
	// check(err)
	return &pb.ReadResponse{}, nil
}

func (server *StorageServer) Write(ctx context.Context, request *pb.WriteRequest) (*pb.WriteResponse, error) {
	filename := request.FileName
	timestamp := request.Timestamp
	val := request.Value

	file := filename + ".json"

	data := Data{
		seconds:  timestamp.Seconds,
		nseconds: timestamp.Nanos,
		value:    val}

	// serialized values from a struct and write them to a file in pretty JSON format
	b, err := json.MarshalIndent(data, "", "	")
	check(err)
	fmt.Println(data)
	fmt.Println(string(b))

	// write to file
	writeResult := ioutil.WriteFile(file, b, 0644)
	check(writeResult)

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
