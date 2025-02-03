# Aurum - Cryptocurrency Wallet Tracker

A high-performance cryptocurrency wallet tracking API with real-time WebSocket support, built in Go using the Fiber framework.

## Features

- 🚀 Real-time wallet balance tracking via WebSocket
- 💰 Multiple blockchain network support
- 📊 Live cryptocurrency price updates
- ⚡ High-performance Go implementation
- 🔒 Secure WebSocket connections with ping/pong health checks

## Prerequisites

- Go 1.23+
- Etherscan API key
- Git (for version tracking)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/xvht/Aurum.git
cd Aurum
```

2. Set up environment:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Run the server:
```bash
go mod download
go run main.go
```

## Build Instructions

### Using Make

Build for your current platform:
```bash
make build
```

Build for specific architectures:
```bash
make build_amd64    # For x86_64 systems
make build_arm64    # For ARM64 systems
make build_windows  # For Windows systems
```

Build for all supported platforms:
```bash
make build_all
```

Debug build:
```bash
make debug
```

Clean build artifacts:
```bash
make clean
```

### Using Docker

1. Build and run using Docker Compose:
```bash
docker-compose up --build
```

2. For production deployment:
```bash
docker-compose -f docker-compose.yml up -d
```

The service will be available at:
- API: `http://localhost:8621`
- WebSocket: `ws://localhost:8621/v1/ws/`

### Architecture Support

Pre-compiled binaries are available for:
- AMD64 (x86_64)
- ARM64
- x86 (386)
- ARM v7
- PPC64/PPC64LE
- s390x
- MIPS64
- Windows (AMD64, 386)

## Deployment

### Docker Environment

The application runs behind Nginx in Docker with the following setup:
- Backend service with automatic restart
- Nginx reverse proxy on port 8621
- Environment configuration via `.env` file
- Health checks and container dependencies handled automatically

## API Documentation

### HTTP Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/health` | GET | Health check endpoint |
| `/v1/version` | GET | Get version information |

### WebSocket Endpoints

#### Query Endpoint: `/v1/ws/query`

Request format:
```json
{
    "queryId": "unique-id",
    "address": "0x...",
    "chain": "ETH"
}
```

Response format:
```json
{
    "queryId": "unique-id",
    "error": false,
    "code": 200,
    "data": {
        "tokenSymbol": "ETH",
        "balance": 1.234,
        "usdValue": 2000.00,
        "lastUpdate": 1677721600000
    }
}
```

#### Prices Endpoint: `/v1/ws/prices`

Provides real-time cryptocurrency price updates.

## Error Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Invalid request |
| 404 | Not found |
| 500 | Server error |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8000 |
| ETHERSCAN_API_KEY | Etherscan API key | - |

## Development

This project uses:
- [Fiber](https://gofiber.io/) for HTTP and WebSocket handling
- [Logrus](https://github.com/sirupsen/logrus) for structured logging
- [godotenv](https://github.com/joho/godotenv) for environment management

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.