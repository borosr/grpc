package main

import (
	"context"
	"fmt"
	"log"
	"log/syslog"
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
	tagname  = "api_log"
	// fluentd  = "127.0.0.1:24224"
	fluentd = "127.0.0.1:5140"
	localD  = "localhost:3333"
)

var storedCount = int32(0)

func main() {
	// syslog.New(syslog.LOG_NOTICE, "api_log")

	// if logwriter, e := syslog.Dial("tcp", localD, syslog.LOG_NOTICE, "api_log"); e == nil {
	if logwriter, e := syslog.Dial("udp", fluentd, syslog.LOG_NOTICE, "api_log"); e == nil {
		log.SetOutput(logwriter)
	} else {
		defer logwriter.Close()
		log.Fatal("Syslog init error!")
	}

	// logger, err := fluent.New(fluent.Config{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer logger.Close()

	// data := map[string]string{
	// 	"foo":  "bar",
	// 	"hoge": "hoge",
	// }

	// if err := logger.Post(tagname, data); err != nil {
	// panic(err)
	// }

	log.Print("Hello log")

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
	message := "First calling\n"
	message += sendToProvider(address1)

	message += "\nSecond calling\n"

	log.Println("Second calling")
	message += sendToProvider(address2)

	w.Write([]byte(message))
}

func sendToProvider(uri string) string {
	start := time.Now()
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := counter.NewCountingClient(conn)
	reply, _ := c.Increment(context.Background(), &counter.CounterRequest{})
	storedCount = reply.Count
	log.Printf("Request duration: %s", time.Since(start))
	return fmt.Sprintf("Count is: %d", reply.Count)
}
