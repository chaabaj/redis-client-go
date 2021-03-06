package main

import (
	"os"

	"github.com/go-redis-client/pkg/client"
)

func main() {
	_, err := client.Connect("localhost", 6379)
	if err != nil {
		println("Cannot connect to redis")
		os.Exit(1)
	}
	println("Connected to redis")
}
