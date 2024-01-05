package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/sebito91/bootdotdev/go/bloggy/api"
)

func main() {
	log.Println("Welcome to Bloggy!")

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	portVal := os.Getenv("PORT")
	port, err := strconv.Atoi(portVal)
	if err != nil {
		panic(fmt.Errorf("could not find value for PORT in env file, please check: %s", err))
	}

	concurrency := flag.Int("concurrency", 10, "number of concurrent feeds to consume")
	sleepInterval := flag.Duration("sleepInterval", time.Minute, "how to to sleep between polls for feeds")

	apiCfg, err := api.GetAPI(*concurrency, *sleepInterval)
	if err != nil {
		panic(err)
	}

	corsMux := api.MiddlewareCors(apiCfg.Router)

	// create the server struct
	server := &http.Server{
		Addr:              fmt.Sprintf("localhost:%d", port),
		Handler:           corsMux,
		ReadHeaderTimeout: time.Second,
	}

	go apiCfg.StartScraping()

	log.Printf("starting bloggy listener on %s...\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	log.Printf("closing bloggy")
}
