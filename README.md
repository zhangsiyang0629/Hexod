# Hexod - Uniswap Market Making Bot

## Description

Hexod is an automated market making bot for Uniswap that compares prices between centralized exchanges (CEX) and Uniswap, executes swaps on Uniswap, and hedges positions on CEX to capture arbitrage opportunities. The project is currently under development.

## Features

- Real-time price monitoring across CEX and Uniswap
- Automated swap execution on Uniswap V3 and V4
- Hedging trades on centralized exchanges
- WebSocket integration for live data feeds
- Configurable trading parameters

## Architecture

The project is structured as follows:

- `config/`: Configuration management
- `constant/`: Constants for exchanges and WebSocket settings
- `router/`: Routing logic for trades
- `types/`: Data type definitions
- `uniswapv3/`: Uniswap V3 client implementation
- `uniswapv4/`: Uniswap V4 client implementation
- `utils/`: Utility functions for math, time, and key generation
- `ws/`: WebSocket client for real-time data

## Installation

1. Ensure you have Go installed (version 1.19 or later)
2. Clone the repository
3. Run `go mod tidy` to install dependencies

## Configuration

Edit `src/config/config.go` to set your API keys, trading parameters, and exchange endpoints.

## Usage

1. Configure your settings in `config.go`
2. Run the bot: `go run main.go` (assuming there's a main.go file)
3. Monitor logs for trade execution

## Development

This project is not yet complete. Planned features include:

- Complete implementation of arbitrage logic
- Risk management modules
- Backtesting capabilities
- Comprehensive error handling

## License

This project is licensed under the MIT License.