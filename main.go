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
	go playWithMetric()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	log.Printf("Recv from %q\n", r.RemoteAddr)
	fmt.Fprintf(w, "Hello DogData!\n")
	fmt.Fprintf(w, "URL.Path=%q, count=%v\n", r.URL.Path, count)
}

func playWithMetric() {
	dogstatsdClient, err := statsd.New("127.0.0.1:8125")

	if err != nil {
		log.Fatal(err)
	}
	for {
		dogstatsdClient.SimpleEvent("An error occurred", "Error message")
		time.Sleep(10 * time.Second)
	}
}
