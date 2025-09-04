package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/websocket"

	"twaporacle/internal/eth"
	"twaporacle/internal/uniswap"

	"github.com/joho/godotenv"
)

type Response struct {
	SpotPrice string `json:"spot_price"`
	TwapPrice string `json:"twap_price"`
	Window    uint32 `json:"window"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(backend bind.ContractBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		window := uint32(900)
		if winParam := r.URL.Query().Get("window"); winParam != "" {
			if parsed, err := strconv.ParseUint(winParam, 10, 32); err == nil {
				if parsed >= 30 && parsed <= 86400 {
					window = uint32(parsed)
				}
			}
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		ctx := context.Background()
		for {
			spotPrice, err := uniswap.SpotPriceWithStruct(ctx, backend)
			if err != nil || spotPrice == nil {
				log.Printf("Error fetching spot price: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			twapPrice, err := uniswap.GetTWAPPrice(ctx, backend, window)
			if err != nil || twapPrice == nil {
				log.Printf("Error fetching TWAP price: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			resp := Response{
				SpotPrice: formatBigFloat(spotPrice),
				TwapPrice: formatBigFloat(twapPrice),
				Window:    window,
			}

			data, _ := json.Marshal(resp)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Write error: %v", err)
				break
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func formatBigFloat(f *big.Float) string {
	return f.Text('f', 18)
}

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	RPC_URL := os.Getenv("API_KEY")

	ethClient, err := eth.NewClient(ctx, RPC_URL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum: %v", err)
	}
	defer ethClient.Close()

	fmt.Println("✅ Connected to Ethereum")

	http.HandleFunc("/ws", wsHandler(ethClient.Client))

	fmt.Println(" WebSocket server running at ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
