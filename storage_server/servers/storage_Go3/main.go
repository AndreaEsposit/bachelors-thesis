package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
)

func main() {
	ExampleStorageServer(8082)
}

func ExampleStorageServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	gorumsSrv := pb.NewGorumsServer()
	srv := storageSrv{state: &pb.State{}}
	gorumsSrv.RegisterStorageServer(&srv)
	gorumsSrv.Serve(lis)
}

type storageSrv struct {
	mut   sync.Mutex
	state *pb.State
}

func (srv *storageSrv) Read(_ context.Context, req *pb.ReadRequest, ret func(*pb.State, error)) {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	fmt.Println("Got Read()")
	ret(srv.state, nil)
}

func (srv *storageSrv) Write(_ context.Context, req *pb.State, ret func(*pb.WriteResponse, error)) {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.state.Timestamp < req.Timestamp {
		srv.state = req
		fmt.Println("Got Write(", req.Value, ")")
		ret(&pb.WriteResponse{New: true}, nil)
		return
	}
	ret(&pb.WriteResponse{New: false}, nil)
}
