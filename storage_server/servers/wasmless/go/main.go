package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// IP is used to choose the IP of the server
const IP = "152.94.162.17:50051" //152.94.162.17:50051 //bbchain7

type StorageServer struct {
	port string
	pb.UnimplementedStorageServer
	mu sync.Mutex
}

type Data struct {
	Seconds  int64  `json:"seconds"`
	Nseconds int32  `json:"nseconds"`
	Value    string `json:"value"`
}

// NewStorageServer initializes an EchoServer
func NewStorageServer() *StorageServer {
	return &StorageServer{
		port: IP, // bbchain21
	}
}

func (server *StorageServer) Read(ctx context.Context, request *pb.ReadRequest) (*pb.ReadResponse, error) {
	filename := request.FileName
	filepath := "./data/" + filename + ".json"

	// defining a struct instance
	var data Data

	server.mu.Lock() // acquire lock
	//f, err := os.Open(file)
	content, err := os.ReadFile(filepath)

	var response = pb.ReadResponse{}
	if err != nil { //os.IsNotExist(err)
		timestamp := timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		}

		// return response
		response.Value = ""
		response.Timestamp = &timestamp
		response.Ok = 0

	} else {
		// decoding data struct
		// from json format
		if e := json.Unmarshal(content, &data); e != nil {
			log.Fatalln("Failed to parse message: ", err)
		}

		timestamp := timestamppb.Timestamp{
			Seconds: data.Seconds,
			Nanos:   data.Nseconds,
		}

		// return response
		response.Value = data.Value
		response.Timestamp = &timestamp
		response.Ok = 1

	}
	server.mu.Unlock() // release lock since error
	return &response, nil
}

func (server *StorageServer) Write(ctx context.Context, request *pb.WriteRequest) (*pb.WriteResponse, error) {
	filen := request.FileName
	timestamp := request.Timestamp
	val := request.Value

	filename := filen + ".json"

	data := Data{
		Seconds:  timestamp.Seconds,
		Nseconds: timestamp.Nanos,
		Value:    val,
	}

	// encode as json in pretty format
	b, err := json.MarshalIndent(data, "", "	")
	check(err)

	server.mu.Lock() // acquire lock
	// write to file
	file, err := os.Open("./data/" + filename)
	check(err)
	file.Write(b)
	file.Sync()
	file.Close()

	//result := os.WriteFile("./data/"+file, b, 0644)
	server.mu.Unlock() // realease lock

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
