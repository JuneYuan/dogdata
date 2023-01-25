package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

var (
	mu    sync.Mutex
	count int
)

func main() {

	go serveStatsd()

	http.HandleFunc("/report", helloReportHandler)
	http.HandleFunc("/query", helloQueryHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func playWithEvent() {
	dogstatsdClient, err := statsd.New("127.0.0.1:8125")

	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		dogstatsdClient.SimpleEvent("An error occurred", "Error message")
		time.Sleep(2 * time.Second)
	}
}

func playWithMetric() {
	client, err := statsd.New("127.0.0.1:8125",
		statsd.WithTags([]string{"env:prod", "service:dogdata"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 20; i++ {
		client.Count("juneyuan.sample", 1, []string{"todayis:great"}, 1)
		fmt.Printf("sent %v value!\n", i+1)
		time.Sleep(100 * time.Millisecond)
	}
	client.Close()
}
