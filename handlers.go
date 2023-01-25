package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func helloReportHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	log.Printf("Recv from %q\n", r.RemoteAddr)
	fmt.Fprintf(w, "Hello DogData!\n")
	fmt.Fprintf(w, "URL.Path=%q, count=%v\n", r.URL.Path, count)
}

func helloQueryHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	log.Printf("Recv from %q\n", r.RemoteAddr)
	fmt.Fprintf(w, "Hello DogData!\n")
	fmt.Fprintf(w, "URL.Path=%q, count=%v\n", r.URL.Path, count)
}

func metricHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "method=%v\n", r.Method)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "io.RealAll(): %v", err)
	}
	fmt.Fprintf(w, "recv metric: %q\n", string(body))
}
