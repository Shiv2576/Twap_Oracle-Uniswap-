package uniswapWeth

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var PoolAddress = common.HexToAddress("0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8")

// Just the slot0 function ABI
const SLOT0_ABI = `[{
	"inputs": [],
	"name": "slot0",
	"outputs": [
		{"internalType": "uint160", "name": "sqrtPriceX96", "type": "uint160"},
		{"internalType": "int24", "name": "tick", "type": "int24"},
		{"internalType": "uint16", "name": "observationIndex", "type": "uint16"},
		{"internalType": "uint16", "name": "observationCardinality", "type": "uint16"},
		{"internalType": "uint16", "name": "observationCardinalityNext", "type": "uint16"},
		{"internalType": "uint8", "name": "feeProtocol", "type": "uint8"},
		{"internalType": "bool", "name": "unlocked", "type": "bool"}
	],
	"stateMutability": "view",
	"type": "function"
}]`

type Slot0 struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

func SpotPriceWithStruct(ctx context.Context, backend bind.ContractBackend) (*big.Float, error) {

	contractABI, err := abi.JSON(strings.NewReader(SLOT0_ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := contractABI.Pack("slot0")
	if err != nil {
		return nil, fmt.Errorf("failed to pack slot0 call: %w", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &PoolAddress,
		Data: data,
	}

	res, err := backend.CallContract(ctx, callMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("eth_call failed: %w", err)
	}

	var slot0Result Slot0
	err = contractABI.UnpackIntoInterface(&slot0Result, "slot0", res)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack slot0 result: %w", err)
	}

	rawPrice := new(big.Float).SetInt(slot0Result.SqrtPriceX96)
	rawPrice.Mul(rawPrice, rawPrice)
	scale := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 192)) // safer than math.Pow
	rawPrice.Quo(rawPrice, scale)

	usdPrice := NormalizePrice(rawPrice, 6, 18)

	return usdPrice, nil
}
