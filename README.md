# AIVM Chain

A working proof-of-concept implementation of the [AIVM Whitepaper](https://docs.chaingpt.org/ai-tools-and-applications/aivm-blockchain-whitepaper) — an AI-native Layer-1 blockchain built on Cosmos SDK.

## What this demonstrates

### 1. AI Smart Contracts (Starlark Engine)
Python-like smart contracts with AI built-ins running in a deterministic sandbox:
```python
def main(args):
    exec_id = request_ai("risk-model", args["input"])
    record = get_execution(exec_id)
    return "Result: " + record["status"]
```
- `request_ai(model_id, input)` — triggers on-chain or off-chain AI inference
- `get_execution(exec_id)` — retrieves execution status and cryptographic proof

### 2. Dual-Path Execution Engine
- **ON_CHAIN**: Simple models run transparently in consensus
- **OFF_CHAIN**: Complex models (LLMs) routed to specialized nodes with proof submission
- Automatic routing based on model type and input complexity

### 3. Cryptographic Proof Verification
- SHA-256 execution proofs anchored on-chain
- Tamper-evident output hashing
- Manual verification confirmed: `PASS ✓`

### 4. CLI Toolkit (`aivmd ai ...`)
- `register-model` — define AI capabilities on-chain
- `execute-contract` — run Starlark AI contracts
- `submit-proof` — validator proof submission bridge

### 5. gRPC Infrastructure
Custom protobuf definitions for high-performance chain ↔ AI worker communication.

## Architecture
```
┌─────────────────────────────────────┐
│  Cosmos SDK Layer-1 Chain (Go)      │
│  ├── x/aimodule (AI Module)         │
│  │   ├── Starlark Contract Engine   │
│  │   ├── Dual-Path Router           │
│  │   ├── Model Registry             │
│  │   └── Proof Verification         │
│  └── CometBFT Consensus             │
└────────────────┬────────────────────┘
                 │ gRPC
┌────────────────▼────────────────────┐
│  Python AI Execution Engine         │
│  ├── FastAPI Node (port 8000)       │
│  ├── OnChain Executor (scikit-learn)│
│  ├── OffChain Executor              │
│  ├── Proof Generator/Verifier       │
│  └── Chain Client (WebSocket)       │
└─────────────────────────────────────┘
```

## Running the demo
```bash
# Terminal 1 — start AI execution node
cd python && source venv/bin/activate
python3 -m uvicorn aivm_engine.node:app --port 8000

# Terminal 2 — run full end-to-end demo
cd python && source venv/bin/activate
python3 demo_full.py
```

## Tech stack
- Cosmos SDK v0.50 + CometBFT
- Go 1.24 + Starlark (go.starlark.net)
- Python 3.12 + FastAPI + scikit-learn
- gRPC + Protocol Buffers
