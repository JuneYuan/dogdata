package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	mu    sync.Mutex
	count int
)

func main1() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8032"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func main() {

	go serveStatsd()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/report", helloReportHandler)
	http.HandleFunc("/query", queryHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8032", nil))
	// log.Fatal(http.ListenAndServe("127.0.0.1:8032", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Docker! <3")
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
