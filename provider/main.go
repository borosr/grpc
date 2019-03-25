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

var storedCount = int32(0)

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
	log.Printf("Counting: %d", storedCount)
	storedCount++
	return &counter.CounterReply{
		Count: storedCount,
	}, nil
}
