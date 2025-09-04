package api

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"

	"twaporacle/internal/uniswap"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Response struct {
	SpotPrice string `json:"spot_price"`
	TwapPrice string `json:"twap_price"`
	Window    int    `json:"window"`
}

func TwapHandler(client *ethclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		var twapWindow uint32 = 900

		if windowParam := r.URL.Query().Get("window"); windowParam != "" {
			if parsed, err := strconv.ParseUint(windowParam, 10, 32); err == nil {
				if parsed >= 30 && parsed <= 86400 {
					twapWindow = uint32(parsed)
				}
			}
		}

		spotPrice, err := uniswap.SpotPriceWithStruct(ctx, client)
		if err != nil || spotPrice == nil {
			http.Error(w, "Failed to fetch spot price: "+err.Error(), http.StatusInternalServerError)
			return
		}

		type TwapResult struct {
			price *big.Float
			err   error
		}
		twapChan := make(chan TwapResult, 1)

		go func() {
			price, err := uniswap.GetTWAPPrice(ctx, client, twapWindow)
			twapChan <- TwapResult{price: price, err: err}
		}()

		twapResult := <-twapChan
		if twapResult.err != nil || twapResult.price == nil {
			http.Error(w, "Failed to fetch TWAP price", http.StatusInternalServerError)
			return
		}

		resp := Response{
			SpotPrice: spotPrice.Text('f', 18),
			TwapPrice: twapResult.price.Text('f', 18),
			Window:    int(twapWindow),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
