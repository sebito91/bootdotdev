package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/sebito91/bootdotdev/go/bloggy/api"
)

func main() {
	fmt.Println("Welcome to Bloggy!")

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	portVal := os.Getenv("PORT")
	port, err := strconv.Atoi(portVal)
	if err != nil {
		panic(fmt.Errorf("could not find value for PORT in env file, please check: %s", err))
	}

	fmt.Printf("we're using port: %d\n", port)

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

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
