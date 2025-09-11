package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"twaporacle/internal/api"
	"twaporacle/internal/eth"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	RPC_URL := os.Getenv("API_KEY")
	if RPC_URL == "" {
		log.Fatal("API_KEY not set in environment")
	}

	ethClient, err := eth.NewClient(ctx, RPC_URL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum: %v", err)
	}
	defer ethClient.Close()

	fmt.Println("Connected to Ethereum")

	http.HandleFunc("/ws", api.WsHandler(ethClient.Client))

	fmt.Println("WebSocket server running at ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
