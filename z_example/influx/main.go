package main

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
)

var (
	bucket = "gospel"
	org    = "DogData"
	token  = "TSmwhXA3hHFZy8i_xAevEDH0iNKfVm4YEw9Wu_mmrCjOcwGUIgNV4rX5VpXA5FLh8nN9zzUSKdyFpLbyOQOBuA=="
	// Store the URL of your InfluxDB instance
	url = "http://localhost:8086"
)

func main() {
	//writeExample()
	readExample()
}

func writeExample() {
	// Create new client with default option for server url authenticate by token
	client := influxdb2.NewClient(url, token)
	// User blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)
	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	// Write point immediately
	writeAPI.WritePoint(context.Background(), p)
	// Ensures background processes finishes
	client.Close()
}

func readExample() {
	// Create client
	client := influxdb2.NewClient(url, token)
	// Get query client
	queryAPI := client.QueryAPI(org)
	// Get QueryTableResult
	result, err := queryAPI.Query(context.Background(), `from(bucket:"gospel")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		// Check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	// Ensures background processes finishes
	client.Close()
}
