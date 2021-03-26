/*
Custom storage-cient (for benchmarking purposes).
Edit IPs with the IPs of the server you want to connect to

Program can be runned like this: go run main.go numberOfRequests mode(r/read/Read/READ or w/Write/write/WRITE)
@author: Andrea Esposito
*/
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

// IPs is used to specify the IPs that we will connect to
var IPs []string

// 10 bytes
const con10b = "Testing!!!"

// 1kb
const con1kb = "Bruce Wayne was born to wealthy physician Dr. Thomas Wayne and his wife Martha, who were themselves members of the prestigious Wayne and Kane families of Gotham City, respectively. When he was three, Bruce's mother Martha was expecting a second child to be named Thomas Wayne, Jr. However, because of her intent to found a school for the underprivileged in Gotham, she was targeted by the manipulative Court of Owls, who arranged for her to have a car accident. She and Bruce survived, but the accident forced Martha into premature labor, and the baby was lost. While on vacation to forget about these events, the Wayne Family butler, Jarvis Pennyworth was killed by one of the Court of Owls' Talons. A letter he'd written to his son Alfred, warning him away from the beleaguered Wayne family, was never delivered. As such, Alfred - who had been an actor at the Globe Theatre at the time and a military medic before that, traveled to Gotham City to take up his father's place, serving the Waynes....."

// 1Mb
var con1Mb = strings.Repeat(con1kb, 100)

var nRequests = 0

func main() {
	//IPs = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054"}
	IPs = []string{"152.94.162.17:50051"} // 152.94.162.12 = bbchain2

	clients := map[int]pb.StorageClient{}

	// creates connections to each server
	for i, ip := range IPs {
		conn, err := grpc.Dial(ip, grpc.WithInsecure())
		check(err)
		fmt.Printf("Connected to: %v\n", ip)
		clients[i] = pb.NewStorageClient(conn)
	}

	// get number of requests
	var err error
	nRequests, err = strconv.Atoi(os.Args[1])
	check(err)

	mode := os.Args[2]

	var latencies []time.Duration
	var doneTimes []int64 //Unix format
	for nRequests != 0 {
		if mode == "w" || strings.ToLower(mode) == "write" {
			timep, err := ptypes.TimestampProto(time.Now())
			check(err)
			message := pb.WriteRequest{FileName: "test", Value: con10b, Timestamp: timep}

			// run requests to all servers specified by IPs
			mWrite(clients, &message, &latencies, &doneTimes)

		} else if mode == "r" || strings.ToLower(mode) == "read" {
			message := pb.ReadRequest{FileName: "test"}
			// run requests to all servers specified by IPs
			mRead(clients, &message, &latencies, &doneTimes)
		}

	}
	// wait before you write to file
	time.Sleep(10 * time.Second)

	file, err := os.Create("result-client" + os.Args[3] + ".csv")
	check(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Latency(MicroSeconds)", "Time(UnixFormat)"})
	checkError("cannot write to file", err)
	for i, value := range latencies {
		s := []string{strconv.Itoa(int(value.Microseconds())), strconv.Itoa(int(doneTimes[i]))}
		err := writer.Write(s)
		checkError("cannot write to file", err)
	}

}

// panic if error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// fatal error making file
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

// mesures latency of a request
func measureTime(latencies *[]time.Duration) func() {
	start := time.Now()
	return func() {
		*latencies = append(*latencies, time.Since(start))
	}
}

func mWrite(clients map[int]pb.StorageClient, message *pb.WriteRequest, latencies *[]time.Duration, times *[]int64) {
	defer measureTime(latencies)()
	for _, client := range clients {
		wg.Add(1)
		go singleWrite(client, message)

	}
	wg.Wait()
	// -1 total requests
	nRequests--
	*times = append(*times, time.Now().Unix())
}

func singleWrite(client pb.StorageClient, message *pb.WriteRequest) {
	_, err := client.Write(context.Background(), message)
	check(err)
	// done current request
	wg.Done()
}

func mRead(clients map[int]pb.StorageClient, message *pb.ReadRequest, latencies *[]time.Duration, times *[]int64) {
	defer measureTime(latencies)()
	for _, client := range clients {
		wg.Add(1)
		go singleRead(client, message)

	}
	wg.Wait()
	// -1 total requests
	nRequests--
	*times = append(*times, time.Now().Unix())
}

func singleRead(client pb.StorageClient, message *pb.ReadRequest) {
	res, _ := client.Read(context.Background(), message)
	if res.GetOk() == 0 {
		err := errors.New("file is not present in one of the servers")
		panic(err)
	}
	// done current request
	wg.Done()
}
