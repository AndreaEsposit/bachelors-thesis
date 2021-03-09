package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

//TODO: Add a re-connection feature if connection is lost

var wg sync.WaitGroup

// IPs is used to specify the IPs that we will connect to
var IPs []string

func main() {
	IPs = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054"}

	clients := map[int]pb.StorageClient{}

	for i, ip := range IPs {
		conn, err := grpc.Dial(ip, grpc.WithInsecure())
		check(err)
		fmt.Printf("Connected to: %v\n", ip)
		clients[i] = pb.NewStorageClient(conn)
	}

	reader := bufio.NewReader(os.Stdout)

	choice := 1

	for choice > 0 {
		fmt.Println("Exit/exit/e to exit this program")
		choice = chooseMode(reader)

		if choice == 1 {
			exit := 0
			for exit == 0 {
				exit, clients = multiRead(reader, clients)
			}
		} else if choice == 2 {
			exit := 0
			for exit == 0 {
				exit, clients = multiWrite(reader, clients)
			}
		}

	}
}

// panic if error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// chooseMode handles the choice of read or write
func chooseMode(reader *bufio.Reader) (choice int) {
	for {
		fmt.Print("Type r/w if you want to read or write data: ")
		command, _ := reader.ReadString('\n')

		command = strings.Replace(command, "\n", "", -1)

		if command == "exit" || command == "Exit" || command == "e" {
			choice = 0
			break
		} else if command == "r" {
			choice = 1
			break
		} else if command == "w" {
			choice = 2
			break
		} else {
			println("Choose either r/w or exit the program")
			continue
		}
	}
	return choice
}

func multiWrite(reader *bufio.Reader, clients map[int]pb.StorageClient) (choice int, activeClients map[int]pb.StorageClient) {
	// Write
	var fName, value string

	// check if name of file has been given
	for {
		fmt.Print("Name of storage (file name): ")
		fName, _ = reader.ReadString('\n')

		fName = strings.Replace(fName, "\n", "", -1)

		if fName == "" {
			continue
		} else {
			break
		}
	}

	// check if a message is not empty
	for {
		fmt.Print("Message to store: ")
		value, _ = reader.ReadString('\n')

		value = strings.Replace(value, "\n", "", -1)

		if value == "" {
			continue
		} else {
			break
		}
	}

	// define timestamp for now-time
	timep, err := ptypes.TimestampProto(time.Now())
	check(err)
	message := pb.WriteRequest{FileName: fName, Value: value, Timestamp: timep}

	//fmt.Printf("We are sending this message: %v\n", message)

	var returnMessages []*pb.WriteResponse
	var failureCalls []int
	for index, client := range clients {
		wg.Add(1)
		c := make(chan wasmRes)
		go singleWrite(client, &message, c, index)
		// remove broken connection
		response := <-c
		if response.err != nil {
			delete(clients, response.i)
			failureCalls = append(failureCalls, response.i)
			// if no more connections are available, panic
			if len(clients) == 0 {
				check(err)
			}
		} else {
			returnMessages = append(returnMessages, response.pb.(*pb.WriteResponse))
		}

	}
	// wait until all writes have been sent and recived
	wg.Wait()

	// check how many succeded
	var successes int
	for _, rMessage := range returnMessages {
		if rMessage.GetOk() == 1 {
			successes++
		}
	}

	// show number of successes
	fmt.Printf("Data stored successfully on %v servers\n", successes)

	// show what connaction failed
	if len(failureCalls) != 0 {
		for _, index := range failureCalls {
			fmt.Printf("Call to IP: %s has failed, connection lost\n", IPs[index])
		}
	}

	fmt.Print("\n'Back/back/b' to go back to mode selection:  ")
	command, _ := reader.ReadString('\n')

	command = strings.Replace(command, "\n", "", -1)
	fmt.Println("")

	if command == "back" || command == "Bxit" || command == "b" {
		return 1, clients
	}
	return 0, clients
}

// singleWrite function will be run by goroutines. It is the actual gRPC write Call
func singleWrite(client pb.StorageClient, message *pb.WriteRequest, c chan (wasmRes), i int) {
	response, err := client.Write(context.Background(), message)
	wg.Done()
	c <- wasmRes{pb: response, err: err, i: i}
}

func multiRead(reader *bufio.Reader, clients map[int]pb.StorageClient) (choice int, activeClients map[int]pb.StorageClient) {
	var fName string

	// check that fName is not black
	for {
		fmt.Print("Name of the storage that you wanna read (file name): ")
		fName, _ = reader.ReadString('\n')

		fName = strings.Replace(fName, "\n", "", -1)

		if fName == "" {
			continue
		} else {
			break
		}
	}

	message := pb.ReadRequest{FileName: fName}

	var returnMessages []*pb.ReadResponse
	var lostConns []int
	var failures []int
	var successes int
	for index, client := range clients {
		wg.Add(1)
		c := make(chan wasmRes)
		go singleRead(client, &message, c, index)
		response := <-c
		// remove broken connection
		if response.err != nil {
			delete(clients, response.i)
			lostConns = append(lostConns, response.i)
			// if no more connections are available, panic
			if len(clients) == 0 {
				check(response.err)
			}
		} else {
			// check if it is a success or a failure
			rMessage := response.pb.(*pb.ReadResponse)
			returnMessages = append(returnMessages, rMessage)
			if rMessage.GetOk() == 1 {
				successes++
			} else {
				failures = append(failures, response.i)
			}
		}

	}
	// wait until all writes have been sent and recived
	wg.Wait()

	// show number of successes
	fmt.Printf("Data read successfully on %v servers\n", successes)

	// show failures
	if len(failures) != 0 {
		for _, index := range failures {
			fmt.Printf("File does not exist on server with IP: %s\n", IPs[index])
		}
	}

	// show what connaction failed
	if len(lostConns) != 0 {
		for _, index := range lostConns {
			fmt.Printf("Call to IP: %s has failed, connection lost\n", IPs[index])
		}
	}

	// return newest content
	n := newest(returnMessages)
	time, err := ptypes.Timestamp(n.Timestamp)
	check(err)
	fmt.Printf("This is the newest (%v) version of the storage: %v ", time, n.Value)

	fmt.Print("\n'Back/back/b' to go back to mode selection:  ")
	command, _ := reader.ReadString('\n')

	command = strings.Replace(command, "\n", "", -1)
	fmt.Println("")

	if command == "back" || command == "Bxit" || command == "b" {
		return 1, clients
	}
	return 0, clients
}

// singleRead function will be run by goroutines. It is the actual gRPC read Call
func singleRead(client pb.StorageClient, message *pb.ReadRequest, c chan (wasmRes), i int) {
	response, err := client.Read(context.Background(), message)
	wg.Done()
	c <- wasmRes{pb: response, err: err, i: i}
}

// returns newesr ReadResponse
func newest(arr []*pb.ReadResponse) *pb.ReadResponse {
	var newest = arr[0]
	for i := range arr {
		if arr[i].Timestamp.Seconds >= newest.Timestamp.Seconds {
			if arr[i].Timestamp.Nanos >= newest.Timestamp.Nanos {
				newest = arr[i]
			}
		}
	}
	return newest
}

// wasmRes is used to communicate between goroutines
type wasmRes struct {
	pb  proto.Message
	err error
	i   int
}
