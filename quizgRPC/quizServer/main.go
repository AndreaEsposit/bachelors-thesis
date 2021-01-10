package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/AndreaEsposit/practice/quizgRPC/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	lis, err := net.Listen("tcp", ":8070") //Listening to snnounvrd on the local network
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterQuizServer(grpcServer, NewQuizServer())
	fmt.Printf("Server is running at :8070.\n")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type QuizServer struct {
	signedUpUsers []*pb.User
	questionList  []*pb.Question
	questionChan  chan *pb.Question
	voteCouter    int
	pb.UnimplementedQuizServer
}

func NewQuizServer() *QuizServer {
	return &QuizServer{
		signedUpUsers: make([]*pb.User, 0),
		questionChan:  make(chan *pb.Question, 10),
	}
}

func (qs *QuizServer) Next(stream pb.Quiz_NextServer) error {
	for {
		question, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			return status.Errorf(codes.NotFound, "Couldn't receive question from quiz master")
		}

		qs.questionChan <- question
		qs.questionList = append(qs.questionList, question)
	}
}

func (qs *QuizServer) Signup(user *pb.User, stream pb.Quiz_SignupServer) error {
	qs.signedUpUsers = append(qs.signedUpUsers, user)
	fmt.Println(qs.signedUpUsers)
	for {
		question := <-qs.questionChan
		err := stream.Send(question)
		if err != nil {
			return status.Errorf(codes.NotFound, "Couldn't send question to %s", user.GetUser())
		}
	}
}

func (qs *QuizServer) Vote(ctx context.Context, vote *pb.VoteRequest) (*pb.User, error) {
	//fmt.Println(qs.signedUpUsers)
	qs.voteCouter++
	for qs.voteCouter < len(qs.signedUpUsers) { //Bad practice, should not just keep spinning, can be fixed in another way (maybe context?),
	}
	qs.resetVoteCounter()

	qId := vote.GetQuestionId()
	fmt.Printf("Correct opt: %v , Your vote: %v\n", qs.questionList[qId-1].GetCorrectOption(), vote.GetVote())
	if qs.questionList[qId-1].GetCorrectOption() == vote.GetVote() {
		fmt.Println("YOU WON")
		return vote.GetUser(), nil
	} else {
		return &pb.User{}, nil
	}
}

func (qs *QuizServer) resetVoteCounter() {
	if qs.voteCouter != 0 {
		qs.voteCouter = 0
	}
}
