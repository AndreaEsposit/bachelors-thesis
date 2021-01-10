package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "github.com/AndreaEsposit/practice/quizgRPC/proto"

	"google.golang.org/grpc"
)

func main() {
	var (
		userName   = flag.String("user", "", "the user name")
		quizMaster = flag.Bool("master", false, "set this to run as quiz master")
	)
	flag.Parse()

	if !*quizMaster && *userName == "" { // Did not choose master and did not choose a name (if chose user)
		flag.Usage()
		return
	}

	conn, err := grpc.Dial(":8070", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewQuizClient(conn) // Make client
	if *quizMaster {                 // You are the quizMaster
		quizMaxter(client)
	} else { // You are a player
		user := &pb.User{User: *userName}
		questionStream, err := client.Signup(context.Background(), user) // Register user in the server
		if err != nil {
			log.Fatal(err)
		}
		for {
			question, err := questionStream.Recv()
			if err != nil {
				fmt.Println("Got an error:", err)
			}
			fmt.Printf("New Question: %d\n-- %s", question.GetId(), question.GetQuestionText())
			for i, q := range question.GetAnswerText() {
				fmt.Printf("---- A%d: %s\n", i, q)
			}
			fmt.Print("What's your answer: ")
			var ansNum int32 // Expect answer as an integer
			fmt.Scanln(&ansNum)

			vote := &pb.VoteRequest{
				QuestionId: question.GetId(),
				Vote:       ansNum,
				User:       user,
			}
			winner, err := client.Vote(context.Background(), vote)

			if err != nil {
				fmt.Println("Got an error:", err)
			}

			if winner.User != user.User {
				fmt.Printf("Sorry, you did not win. Try your luck with the next question!\n")
				if winner.User == "" { // Noone guessed right
					fmt.Printf("Nobody won round %d\n", vote.QuestionId)
				} else {
					fmt.Printf("The winner of round %d was %s\n", vote.QuestionId, winner.User)
				}

			} else {
				fmt.Println("Congratulations! You are the winner of this round!")
			}

		}
	}
}

func quizMaxter(client pb.QuizClient) {
	stream, err := client.Next(context.Background())
	if err != nil {
		fmt.Println("Got an error:", err)
	}
	questionTable := []*pb.Question{
		{
			Id:            1,
			QuestionText:  "2 + 2 = ?\n",
			AnswerText:    []string{"2", "4", "5", "22"},
			CorrectOption: 1,
		},
		{
			Id:            2,
			QuestionText:  "Can we go home now?\n",
			AnswerText:    []string{"Not yet", "Soon", "Never", "Tomorrow", "Yes"},
			CorrectOption: 2,
		},
		{
			Id:            3,
			QuestionText:  "What comes before 1?\n",
			AnswerText:    []string{"2", "1", "-1", "700", "0"},
			CorrectOption: 4,
		},
	}

	for _, q := range questionTable {
		fmt.Printf("Sending question %d\n", q.GetId())
		err = stream.Send(q)
		if err != nil {
			fmt.Println("Got an error:", err)
		}
	}
	for {
	} //keep stream open

}
