package uniswapLink

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const OBSERVE_ABI = `[{
	"inputs": [
		{"internalType": "uint32[]", "name": "secondsAgos", "type": "uint32[]"}
	],
	"name": "observe",
	"outputs": [
		{"internalType": "int56[]", "name": "tickCumulatives", "type": "int56[]"},
		{"internalType": "uint160[]", "name": "secondsPerLiquidityCumulativeX128s", "type": "uint160[]"}
	],
	"stateMutability": "view",
	"type": "function"
}]`

type ObserveResult struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}

func CumulativePrice(ctx context.Context, backend bind.ContractBackend, secondsAgo []uint32) (*ObserveResult, error) {
	contractAbi, err := abi.JSON(strings.NewReader(OBSERVE_ABI))

	if err != nil {
		return nil, fmt.Errorf("failed at contractABI: %w", err)
	}

	data, err := contractAbi.Pack("observe", secondsAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to pack Data: %w", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &PoolAddress,
		Data: data,
	}

	res, err := backend.CallContract(ctx, callMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	var result ObserveResult
	err = contractAbi.UnpackIntoInterface(&result, "observe", res)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack observe result: %w", err)
	}

	return &result, nil
}

func CalculateTWAP(ctx context.Context, backend bind.ContractBackend, windowSeconds uint32) (*big.Float, error) {

	secondsAgo := []uint32{0, windowSeconds}

	result, err := CumulativePrice(ctx, backend, secondsAgo)
	if err != nil {
		return nil, fmt.Errorf("failed at CummalativePrice: %w", err)
	}

	if len(result.TickCumulatives) < 2 {
		return nil, fmt.Errorf("insufficient cumulative data")
	}

	tickDiff := new(big.Int).Sub(result.TickCumulatives[0], result.TickCumulatives[1])
	timeDiff := new(big.Int).SetUint64(uint64(windowSeconds))

	avgtick := new(big.Float).SetInt(tickDiff)
	avgtick.Quo(avgtick, new(big.Float).SetInt(timeDiff))

	tickFloat, _ := avgtick.Float64()
	price := calculatePriceFromTick(tickFloat)

	return price, nil

}

func calculatePriceFromTick(tick float64) *big.Float {

	ln1_0001 := 0.00009999500033330835 // ln(1.0001)
	priceFloat := math.Exp(tick * ln1_0001)

	return new(big.Float).SetFloat64(priceFloat)

}

func GetTWAPPrice(ctx context.Context, backend bind.ContractBackend, windowSeconds uint32) (*big.Float, error) {
	twapPrice, err := CalculateTWAP(ctx, backend, windowSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate TWAP: %w", err)
	}

	decimals0 := 18 //LINK
	decimals1 := 18 //ETH

	usdTwapPrice := NormalizePrice(twapPrice, decimals0, decimals1)
	return usdTwapPrice, nil
}

func NormalizePrice(rawPrice *big.Float, decimals0, decimals1 int) *big.Float {
	exp := decimals1 - decimals0

	factor := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil),
	)
	normalized := new(big.Float).Mul(rawPrice, factor)

	return normalized
}
