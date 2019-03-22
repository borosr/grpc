package main

import (
	"context"
	"log"
	"net"

	counter "github.com/borosr/go_grpc"
	"google.golang.org/grpc"
)

const (
	providerPort = ":8086"
)

type server struct{}

func main() {
	log.Println("Provider started")
	lis, err := net.Listen("tcp", providerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	counter.RegisterCountingServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Increment ....
func (s *server) Increment(ctx context.Context, in *counter.CounterRequest) (*counter.CounterReply, error) {
	log.Printf("Counting: %d", in.Count)
	return &counter.CounterReply{
		Count: in.Count + 1,
	}, nil
}
