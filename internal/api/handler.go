package api

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"twaporacle/internal/uniswapAave"
	"twaporacle/internal/uniswapDai"
	"twaporacle/internal/uniswapLink"
	"twaporacle/internal/uniswapOhm"
	"twaporacle/internal/uniswapPepe"
	"twaporacle/internal/uniswapUni"
	"twaporacle/internal/uniswapWbtc"
	"twaporacle/internal/uniswapWeth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/websocket"
)

type PoolPrice struct {
	SpotPrice string `json_string:"spot_price"`
	TwapPrice string `json_string:"twap_price"`
	Slippage  string `json:"slippage_percent,omitempty"`
}

type Response struct {
	Pools  map[string]PoolPrice `json_string:"pools"`
	Window uint32               `json_string:"window"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsHandler(backend bind.ContractBackend) http.HandlerFunc {
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
			pools := make(map[string]PoolPrice)

			// WETH/USDC
			WethSpot, err := uniswapWeth.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching WETH spot price: %v", err)
			} else {
				WethTwap, err := uniswapWeth.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching WETH TWAP price: %v", err)
				} else {
					wethSlippage := CalculateSlippage(WethSpot, WethTwap)
					pools["WETH/USDC"] = PoolPrice{
						SpotPrice: formatBigFloatWETH(WethSpot),
						TwapPrice: formatBigFloatWETH(WethTwap),
						Slippage:  formatBigFloatSlip(wethSlippage),
					}
				}
			}

			// DAI/USDC
			DaiSpot, err := uniswapDai.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching DAI spot price: %v", err)
			} else {
				DaiTwap, err := uniswapDai.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching DAI TWAP price: %v", err)
				} else {
					daiSlippage := CalculateSlippage(DaiSpot, DaiTwap)
					pools["DAI/USDC"] = PoolPrice{
						SpotPrice: formatBigFloatDAI(DaiSpot),
						TwapPrice: formatBigFloatDAI(DaiTwap),
						Slippage:  formatBigFloatSlip(daiSlippage),
					}
				}
			}

			//WBTC/USDC
			WbtcSpot, err := uniswapWbtc.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching WBTC spot price: %v", err)
			} else {
				WbtcTwap, err := uniswapWbtc.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching WBTC TWAP price: %v", err)
				} else {
					wbtcSlippage := CalculateSlippage(WbtcSpot, WbtcTwap)
					pools["WBTC/USDC"] = PoolPrice{
						SpotPrice: formatBigFloatWETH(WbtcSpot),
						TwapPrice: formatBigFloatWETH(WbtcTwap),
						Slippage:  formatBigFloatSlip(wbtcSlippage),
					}
				}
			}

			//LINK/ETH
			LinkSpot, err := uniswapLink.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching YES spot price: %v", err)
			} else {
				LinkTwap, err := uniswapLink.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching YES TWAP price: %v", err)
				} else {
					linkSlippage := CalculateSlippage(LinkSpot, LinkTwap)
					pools["LINK/ETH"] = PoolPrice{
						SpotPrice: formatBigFloatWETH(LinkSpot),
						TwapPrice: formatBigFloatWETH(LinkTwap),
						Slippage:  formatBigFloatSlip(linkSlippage),
					}
				}
			}

			//UNI/ETH
			UniSpot, err := uniswapUni.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching UNI spot price: %v", err)
			} else {
				UniTwap, err := uniswapUni.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching UNI TWAP price: %v", err)
				} else {
					uniSlippage := CalculateSlippage(UniSpot, UniTwap)
					pools["UNI/ETH"] = PoolPrice{
						SpotPrice: formatBigFloatUni(UniSpot),
						TwapPrice: formatBigFloatUni(UniTwap),
						Slippage:  formatBigFloatSlip(uniSlippage),
					}
				}
			}

			//AAVE/ETH
			AaveSpot, err := uniswapAave.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching TRX spot price: %v", err)
			} else {
				AaveTwap, err := uniswapAave.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching TRX TWAP price: %v", err)
				} else {
					AaveSlippage := CalculateSlippage(AaveSpot, AaveTwap)
					pools["AAVE/ETH"] = PoolPrice{
						SpotPrice: formatBigFloatUni(AaveSpot),
						TwapPrice: formatBigFloatUni(AaveTwap),
						Slippage:  formatBigFloatSlip(AaveSlippage),
					}
				}
			}

			//PEPE/ETH
			PepeSpot, err := uniswapPepe.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching TRX spot price: %v", err)
			} else {
				PepeTwap, err := uniswapPepe.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching TRX TWAP price: %v", err)
				} else {
					PepeSlippage := CalculateSlippage(PepeSpot, PepeTwap)
					pools["PEPE/ETH"] = PoolPrice{
						SpotPrice: formatBigFloatPepe(PepeSpot),
						TwapPrice: formatBigFloatPepe(PepeTwap),
						Slippage:  formatBigFloatSlip(PepeSlippage),
					}
				}
			}

			//OHM/ETH
			OhmSpot, err := uniswapOhm.SpotPriceWithStruct(ctx, backend)
			if err != nil {
				log.Printf("Error fetching TRX spot price: %v", err)
			} else {
				OhmTwap, err := uniswapOhm.GetTWAPPrice(ctx, backend, window)
				if err != nil {
					log.Printf("Error fetching TRX TWAP price: %v", err)
				} else {
					OhmSlippage := CalculateSlippage(OhmSpot, OhmTwap)
					pools["OHM/ETH"] = PoolPrice{
						SpotPrice: formatBigFloatPepe(OhmSpot),
						TwapPrice: formatBigFloatPepe(OhmTwap),
						Slippage:  formatBigFloatSlip(OhmSlippage),
					}
				}
			}

			if len(pools) > 0 {
				resp := Response{
					Pools:  pools,
					Window: window,
				}

				data, err := json.Marshal(resp)
				if err != nil {
					log.Printf("JSON marshaling error: %v", err)
				} else {
					if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
						log.Printf("WebSocket write error: %v", err)
						break
					}
				}
			} else {
				log.Printf("No pool data available in this iteration")
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func formatBigFloatWETH(f *big.Float) string {
	return f.Text('f', 8)
}

func formatBigFloatUni(f *big.Float) string {
	return f.Text('f', 8)
}

func formatBigFloatPepe(f *big.Float) string {
	return f.Text('f', 18)
}

func formatBigFloatSlip(f *big.Float) string {
	return f.Text('f', 8)
}

func formatBigFloatDAI(f *big.Float) string {
	factor := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil), // 10^12
	)
	// DIVIDE by 10^12
	adjustedPrice := new(big.Float).Quo(f, factor)

	return adjustedPrice.Text('f', 8)
}

func CalculateSlippage(spotPrice, twapPrice *big.Float) *big.Float {
	// Calculate absolute difference: |spot - twap|
	diff := new(big.Float).Sub(spotPrice, twapPrice)
	absDiff := new(big.Float).Abs(diff)

	// Calculate slippage percentage: (|diff| / twap) * 100
	slippage := new(big.Float).Quo(absDiff, twapPrice)
	slippage.Mul(slippage, big.NewFloat(100))

	return slippage
}
