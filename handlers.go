package main

import (
	"fmt"
	"github.com/kr/pretty"
	"io"
	"log"
	"net/http"
	"sideproject/dogdata/datastore/influx"
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

func queryHandler(w http.ResponseWriter, r *http.Request) {
	values, err := influx.NewWrapClient(url, token).Query()
	if err != nil {
		fmt.Fprintf(w, "Query(): %v", err)
	}

	fmt.Fprintf(w, "Query(): %v", pretty.Sprint(values))
}

func metricHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "method=%v\n", r.Method)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "io.RealAll(): %v", err)
	}
	fmt.Fprintf(w, "recv metric: %q\n", string(body))
}
