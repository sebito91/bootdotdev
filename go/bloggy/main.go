package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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
}
