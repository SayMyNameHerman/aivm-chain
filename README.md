# AIVM Chain

An AI-native Layer-1 blockchain built on Cosmos SDK, implementing the core architectural concepts from the [AIVM Whitepaper](https://docs.chaingpt.org/ai-tools-and-applications/aivm-blockchain-whitepaper).

## Overview

AIVM Chain demonstrates a working proof-of-concept of decentralized AI execution infrastructure, with a custom AI module implementing dual-path execution — the core innovation described in the AIVM whitepaper.

## Architecture

### Dual-Path Execution Engine
The core innovation: AI workloads are routed based on complexity.

- **ON_CHAIN**: Simple models execute directly on-chain with full transparency and consensus-based validation
- **OFF_CHAIN**: Complex models (LLMs, computer vision) route to specialized execution nodes, with cryptographic proofs submitted back on-chain

### AI Module (`x/aimodule`)
Custom Cosmos SDK module implementing:
- **Model Registry**: Register and manage AI models on-chain with immutable hash verification
- **Execution Router**: Dual-path routing logic based on model type
- **Proof Verification**: SHA-256 based cryptographic proof submission and verification for off-chain executions

## Technical Stack
- Cosmos SDK v0.50
- CometBFT consensus
- Go 1.24

## Running Locally
```bash
# Build
go build -o aivmd ./cmd/minid

# Initialize
./aivmd init mynode --chain-id aivm-1
./aivmd keys add validator --keyring-backend test
./aivmd genesis add-genesis-account validator 10000000000stake --keyring-backend test
./aivmd genesis gentx validator 1000000000stake --chain-id aivm-1 --keyring-backend test
./aivmd genesis collect-gentxs

# Start
./aivmd start --minimum-gas-prices 0stake
```

## Roadmap
- [ ] Python AI execution engine (off-chain node)
- [ ] gRPC bridge between Go chain and Python AI layer
- [ ] Validator specialization (AI, Compute, Data validators)
- [ ] REST API for model registration and execution requests
- [ ] Token economics and staking mechanics
