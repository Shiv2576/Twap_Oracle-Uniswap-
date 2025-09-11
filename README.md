# Twap_Oracle-Uniswap-

# TWAP Oracle

TWAP Oracle is a Go-based WebSocket service that streams both spot prices and time-weighted average prices (TWAP) for the USDC/WETH pair on Uniswap. It helps traders simulate swaps with minimal slippage and detect arbitrage opportunities when spot and TWAP diverge.

## Features

- Real-time spot price streaming
- Multiple TWAP window calculations (5min, 15min, 30min, 1hr)
- WebSocket-based streaming for low latency
- Support for multiple client languages
- Arbitrage opportunity detection

## Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/Shiv2576/Twap_Oracle-Uniswap-.git
cd Twap_Oracle-Uniswap-
go mod tidy
```

## Usage

### Start the WebSocket server:

```bash
go run cmd/server/main.go
```

### Connect to the server:

**Spot price stream:**
```
ws://localhost:8080/ws?window=spot
```

**TWAP streams:**
```
ws://localhost:8080/ws?window=300    # 5 minutes
ws://localhost:8080/ws?window=900    # 15 minutes
ws://localhost:8080/ws?window=1800   # 30 minutes
ws://localhost:8080/ws?window=3600   # 60 minutes
```

## Example Clients

### Go

```go
conn, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/ws?window=300", nil)
```

## API Reference

### WebSocket Endpoints

| Endpoint | Description | Parameters |
|----------|-------------|------------|
| `/ws` | Main WebSocket endpoint | `window`: Time window in seconds or "spot" |

### Supported Windows

- `spot`: Real-time spot prices
- `300`: 5-minute TWAP
- `900`: 15-minute TWAP
- `1800`: 30-minute TWAP
- `3600`: 1-hour TWAP
