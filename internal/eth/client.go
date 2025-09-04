package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	*ethclient.Client
}

func NewClient(ctx context.Context, rpcURL string) (*Client, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("Failed To connect To Ethereum RPC :  %w", err)
	}

	return &Client{client}, err
}
