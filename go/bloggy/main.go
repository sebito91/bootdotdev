package main

import (
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

	apiCfg, err := api.GetAPI()
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

	// TODO: move these to flags
	var concurrency = 10
	var sleepInterval = time.Minute
	go apiCfg.StartScraping(concurrency, sleepInterval)

	log.Printf("starting bloggy listener on %s...\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	log.Printf("closing bloggy")
}
