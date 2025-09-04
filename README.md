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

### Python

```python
import asyncio
import websockets

async def listen():
    url = "ws://localhost:8080/ws?window=300"
    async with websockets.connect(url) as ws:
        while True:
            msg = await ws.recv()
            print(msg)

asyncio.run(listen())
```

### Go

```go
conn, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/ws?window=300", nil)
```

### JavaScript

```javascript
const ws = new WebSocket("ws://localhost:8080/ws?window=300");
ws.on("message", (data) => console.log(data.toString()));
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

## Use Cases

- **Arbitrage Detection**: Compare spot vs TWAP prices to identify profit opportunities
- **Slippage Simulation**: Use TWAP data to estimate swap execution with minimal slippage
- **Market Analysis**: Monitor price trends across different time windows
- **Trading Bots**: Real-time data feed for algorithmic trading strategies

## Requirements

- Go 1.19 or higher
- Internet connection for Uniswap data
- WebSocket-compatible client

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`
