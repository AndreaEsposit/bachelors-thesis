package main

import (
	"context"
	"log"
	"time"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"google.golang.org/grpc"
)

func main() {
	ExampleStorageClient()
}

func ExampleStorageClient() {
	addrs := []string{
		"127.0.0.1:8080",
		"127.0.0.1:8081",
		"127.0.0.1:8082",
	}

	mgr, err := pb.NewManager(pb.WithNodeList(addrs), pb.WithGrpcDialOptions(
		grpc.WithBlock(),
		grpc.WithInsecure(),
	),
		pb.WithDialTimeout(500*time.Millisecond),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Get all all available node ids, 3 nodes
	ids := mgr.NodeIDs()

	// Create a configuration including all nodes
	allNodesConfig, err := mgr.NewConfiguration(ids, nil)
	if err != nil {
		log.Fatalln("error creating read config:", err)
	}
	// Test state
	state := &pb.State{
		Value:     "42",
		Timestamp: time.Now().Unix(),
	}

	// Invoke Write RPC on all nodes in config
	for _, node := range allNodesConfig.Nodes() {
		respons, err := node.Write(context.Background(), state)
		if err != nil {
			log.Fatalln("read rpc returned error:", err)
		} else if !respons.New {
			log.Println("state was not new.")
		}
	}
}
