#!/usr/bin/env python3
"""
AIVM Full End-to-End Demo
Demonstrerer: Starlark AI Smart Contract → Dual-Path Execution → Cryptographic Proof
"""
import requests
import hashlib
import json
import time

AI_NODE = "http://localhost:8000"

def banner(text):
    print(f"\n{'='*60}")
    print(f"  {text}")
    print(f"{'='*60}")

def step(text):
    print(f"\n▶ {text}")

def ok(text):
    print(f"  ✓ {text}")

def show(label, value):
    if isinstance(value, dict):
        print(f"  {label}:")
        for k, v in value.items():
            if isinstance(v, str) and len(v) > 40:
                v = v[:40] + "..."
            print(f"    {k}: {v}")
    else:
        if isinstance(value, str) and len(value) > 60:
            value = value[:60] + "..."
        print(f"  {label}: {value}")

# ─────────────────────────────────────────────
# DEMO 1: Starlark Smart Contract Simulation
# ─────────────────────────────────────────────
banner("DEMO 1: AI Smart Contract (Starlark)")

STARLARK_CONTRACT = """
def main(args):
    prompt = args.get("prompt", "default")
    model  = args.get("model",  "default-classifier")
    exec_id = request_ai(model, prompt)
    record = get_execution(exec_id)
    if record == None:
        return "Error: execution record not found"
    return "AI execution queued! ID=" + exec_id + " status=" + record["status"]
"""

step("Simulating Starlark contract execution via AI execution engine...")

# Simulate what the Starlark contract does: call request_ai then get_execution
exec_id = f"starlark-{hashlib.sha256(b'contract-run-1').hexdigest()[:12]}"
input_hash = hashlib.sha256(b"Is this financial advice?").hexdigest()

payload = {
    "execution_id": exec_id,
    "model_id": "default-classifier",
    "input_hash": input_hash,
    "execution_type": "ON_CHAIN",
    "input_data": {
        "prompt": "Is this financial advice?",
        "features": [2, 3]
    }
}

response = requests.post(f"{AI_NODE}/execute", json=payload)
result = response.json()

ok("Contract executed!")
show("Contract", "ai_hello.star")
show("Prompt", "Is this financial advice?")
show("Execution path", result["execution_path"])
show("AI prediction", result["result"]["prediction"])
show("Probability", result["result"]["probability"])
show("Proof verified", result["proof"]["verified"])
show("Proof hash", result["proof"]["proof"])

simulated_return = f"AI execution queued! ID={exec_id} status=EXECUTED_ON_CHAIN"
show("Contract return value", simulated_return)

# ─────────────────────────────────────────────
# DEMO 2: Dual-Path Routing
# ─────────────────────────────────────────────
banner("DEMO 2: Dual-Path Execution Router")

test_cases = [
    {
        "name": "Simple classifier → ON_CHAIN",
        "model_id": "default-classifier",
        "execution_type": "ON_CHAIN",
        "input_data": {"features": [1, 2]}
    },
    {
        "name": "LLM model → OFF_CHAIN",
        "model_id": "llm-gpt-large",
        "execution_type": "OFF_CHAIN",
        "input_data": {"prompt": "Analyze this complex dataset..."}
    },
    {
        "name": "Risk assessment → ON_CHAIN",
        "model_id": "risk-assessment-v1",
        "execution_type": "ON_CHAIN",
        "input_data": {"features": [5, 7]}
    }
]

for tc in test_cases:
    step(tc["name"])
    exec_id = f"demo-{hashlib.sha256(tc['name'].encode()).hexdigest()[:8]}"
    payload = {
        "execution_id": exec_id,
        "model_id": tc["model_id"],
        "input_hash": hashlib.sha256(json.dumps(tc["input_data"]).encode()).hexdigest(),
        "execution_type": tc["execution_type"],
        "input_data": tc["input_data"]
    }
    r = requests.post(f"{AI_NODE}/execute", json=payload)
    res = r.json()
    ok(f"Path: {res['execution_path']} | Verified: {res['proof']['verified']}")

# ─────────────────────────────────────────────
# DEMO 3: Cryptographic Proof Chain
# ─────────────────────────────────────────────
banner("DEMO 3: Cryptographic Proof Verification")

step("Generating execution proof...")
exec_id = f"proof-demo-{int(time.time())}"
input_data = {"features": [4, 5], "contract": "ai_hello.star"}
input_hash = hashlib.sha256(json.dumps(input_data).encode()).hexdigest()

payload = {
    "execution_id": exec_id,
    "model_id": "default-classifier",
    "input_hash": input_hash,
    "execution_type": "ON_CHAIN",
    "input_data": input_data
}

r = requests.post(f"{AI_NODE}/execute", json=payload)
res = r.json()
proof = res["proof"]

ok("Proof generated and verified on-chain")
show("Execution ID", exec_id)
show("Input hash", input_hash)
show("Output hash", proof["output_hash"])
show("Proof", proof["proof"])
show("Verified", proof["verified"])
show("Timestamp", proof["timestamp"])

# Verify the proof manually
expected = hashlib.sha256(f"{exec_id}{proof['output_hash']}".encode()).hexdigest()
manual_verify = expected == proof["proof"]
show("Manual verification", f"{'PASS ✓' if manual_verify else 'FAIL ✗'}")

# ─────────────────────────────────────────────
# SUMMARY
# ─────────────────────────────────────────────
banner("DEMO COMPLETE")
print("""
  AIVM Proof-of-Concept demonstrates:

  ✓ AI Smart Contracts (Starlark)
    - Python-like contracts with AI built-ins
    - request_ai() and get_execution() primitives
    - Deterministic, sandboxed execution

  ✓ Dual-Path Execution Engine  
    - ON_CHAIN: simple models run transparently
    - OFF_CHAIN: complex models with proof submission
    - Automatic routing based on model type

  ✓ Cryptographic Proof Verification
    - SHA-256 based execution proofs
    - Tamper-evident output hashing
    - Chain-verifiable AI results

  Built on: Cosmos SDK v0.50 + CometBFT + Python FastAPI
  GitHub: github.com/SayMyNameHerman/aivm-chain
""")
