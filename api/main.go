package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	counter "github.com/borosr/go_grpc"
	"google.golang.org/grpc"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	address1 = "localhost:8086"
	address2 = "localhost:8087"
	apiPort  = ":3000"
)

var storedCount = int32(0)

func main() {
	log.Println("API started")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handleCounterGet)

	http.ListenAndServe(apiPort, r)
}

func handleCounterGet(w http.ResponseWriter, r *http.Request) {
	log.Println("First calling")
	sendToProvider(address1)

	log.Println("Second calling")
	message := sendToProvider(address2)

	w.Write(message)
}

func sendToProvider(uri string) []byte {
	start := time.Now()
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := counter.NewCountingClient(conn)
	reply, _ := c.Increment(context.Background(), &counter.CounterRequest{Count: storedCount})
	storedCount = reply.Count
	log.Printf("Request duration: %s", time.Since(start))
	return []byte(fmt.Sprintf("Count is: %d", reply.Count))
}
