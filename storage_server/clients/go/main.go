package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/AndreaEsposit/practice/storage_server/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //152.94.1.100
	check(err)

	client := pb.NewStorageClient(conn)
	reader := bufio.NewReader(os.Stdout)

	fmt.Println("Exit/exit/e' to exit this program")
	choice := chooseMode(reader)

	if choice == 1 {
		exit := 0
		for exit == 0 {
			exit = read(reader, client)
		}
	} else {
		exit := 0
		for exit == 0 {
			exit = write(reader, client)
		}
	}

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func chooseMode(reader *bufio.Reader) (choice int) {
	for {
		fmt.Print("Type r/w if you want to read or write data: ")
		command, _ := reader.ReadString('\n')

		command = strings.Replace(command, "\n", "", -1)

		if command == "exit" || command == "Exit" || command == "e" {
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

func write(reader *bufio.Reader, client pb.StorageClient) (choice int) {
	// Write
	var fName, value string

	for {
		fmt.Print("Name of storage (file name): ")
		fName, _ := reader.ReadString('\n')

		fName = strings.Replace(fName, "\n", "", -1)

		if fName == "" {
			continue
		} else {
			break
		}
	}
	for {
		fmt.Print("Message to store: ")
		value, _ := reader.ReadString('\n')

		value = strings.Replace(value, "\n", "", -1)

		if value == "" {
			continue
		} else {
			break
		}
	}

	timep, err := ptypes.TimestampProto(time.Now())
	message := pb.WriteRequest{FileName: fName, Value: value, Timestamp: timep}

	returnMessage, err := client.Write(context.Background(), &message)
	check(err)

	if returnMessage.GetOk() == 0 {
		fmt.Println("Something went wrong storing your data")

	} else {
		fmt.Println("Data stored successfully")
	}

	fmt.Print("'Exit/exit/e' to go back to mode selection:  ")
	command, _ := reader.ReadString('\n')

	command = strings.Replace(value, "\n", "", -1)

	if command == "exit" || command == "Exit" || command == "e" {
		return 1
	}
	return 0
}

func read(reader *bufio.Reader, client pb.StorageClient) (choice int) {
	// Write
	var fName, value string

	for {
		fmt.Print("Name of the storage that you wanna read (file name): ")
		fName, _ := reader.ReadString('\n')

		fName = strings.Replace(fName, "\n", "", -1)

		if fName == "" {
			continue
		} else {
			break
		}
	}

	message := pb.ReadRequest{FileName: fName}

	returnMessage, err := client.Read(context.Background(), &message)
	check(err)

	fmt.Printf("This is the content that you have recived from the server: %v\n", returnMessage)

	fmt.Print("'Exit/exit/e' to go back to mode selection:  ")
	command, _ := reader.ReadString('\n')

	command = strings.Replace(value, "\n", "", -1)

	if command == "exit" || command == "Exit" || command == "e" {
		return 1
	}
	return 0
}
