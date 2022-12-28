package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("go caching with redis")

	err := godotenv.Load()
	check(err)

	ctx := context.Background()

	redisClient, err := NewRedisClient(ctx)
	check(err)
	defer redisClient.Close()

	postgresClient, err := NewPostgresClient(ctx)
	check(err)
	defer postgresClient.Close()

	server := NewServer(os.Getenv("SERVER_PORT"), redisClient, postgresClient)
	server.start()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
