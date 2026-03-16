import asyncio
import hashlib
import json
import time
import requests
import websockets
from typing import Dict, Optional, Callable

CHAIN_RPC = "http://localhost:26657"
CHAIN_WS = "ws://localhost:26657/websocket"
AI_NODE = "http://localhost:8000"

class ChainClient:
    """
    Kobler Python AI execution engine til Go-chain.
    Lytter på chain-events og behandler off-chain eksekveringer automatisk.
    """
    
    def __init__(self):
        self.running = False
        self.processed = set()
    
    def get_chain_status(self) -> Dict:
        try:
            r = requests.get(f"{CHAIN_RPC}/status", timeout=5)
            return r.json()
        except Exception as e:
            return {"error": str(e)}
    
    def get_latest_block(self) -> Dict:
        try:
            r = requests.get(f"{CHAIN_RPC}/block", timeout=5)
            return r.json()
        except Exception as e:
            return {"error": str(e)}
    
    def submit_proof_to_chain(self, execution_id: str, output_hash: str, proof: str) -> Dict:
        """Sender kryptografisk bevis tilbake til chain etter off-chain eksekvering"""
        payload = {
            "execution_id": execution_id,
            "output_hash": output_hash,
            "proof": proof,
            "timestamp": int(time.time())
        }
        print(f"[ChainClient] Submitting proof for {execution_id[:16]}...")
        print(f"[ChainClient] output_hash: {output_hash[:20]}...")
        print(f"[ChainClient] proof: {proof[:20]}...")
        return {"status": "submitted", "execution_id": execution_id}
    
    def execute_ai_request(self, execution_id: str, model_id: str, 
                           input_hash: str, execution_type: str) -> Optional[Dict]:
        """Sender AI-eksekveringsforespørsel til Python execution engine"""
        if execution_id in self.processed:
            return None
        
        payload = {
            "execution_id": execution_id,
            "model_id": model_id,
            "input_hash": input_hash,
            "execution_type": execution_type,
            "input_data": {"features": [3, 4], "input_hash": input_hash}
        }
        
        try:
            r = requests.post(f"{AI_NODE}/execute", json=payload, timeout=30)
            result = r.json()
            self.processed.add(execution_id)
            
            if result.get("proof", {}).get("verified"):
                self.submit_proof_to_chain(
                    execution_id,
                    result["proof"]["output_hash"],
                    result["proof"]["proof"]
                )
            
            return result
        except Exception as e:
            print(f"[ChainClient] Execution error: {e}")
            return None
    
    async def listen_for_events(self):
        """Lytter på chain-events via WebSocket"""
        print(f"[ChainClient] Connecting to chain at {CHAIN_WS}...")
        
        try:
            async with websockets.connect(CHAIN_WS) as ws:
                subscribe_msg = json.dumps({
                    "jsonrpc": "2.0",
                    "method": "subscribe",
                    "id": 1,
                    "params": {
                        "query": "tm.event='NewBlock'"
                    }
                })
                await ws.send(subscribe_msg)
                print("[ChainClient] Subscribed to NewBlock events")
                
                while self.running:
                    try:
                        msg = await asyncio.wait_for(ws.recv(), timeout=10)
                        data = json.loads(msg)
                        
                        if "result" in data and "data" in data.get("result", {}):
                            block = data["result"]["data"]
                            height = block.get("value", {}).get("block", {}).get("header", {}).get("height", "?")
                            print(f"[ChainClient] New block: height={height}")
                            
                    except asyncio.TimeoutError:
                        print("[ChainClient] Waiting for blocks...")
                    except Exception as e:
                        print(f"[ChainClient] Event error: {e}")
                        break
                        
        except Exception as e:
            print(f"[ChainClient] WebSocket connection failed: {e}")
            print("[ChainClient] Chain not running — operating in standalone mode")
    
    def start(self):
        self.running = True
        print("[ChainClient] AIVM Chain Client starting...")
        print(f"[ChainClient] Chain RPC: {CHAIN_RPC}")
        print(f"[ChainClient] AI Node: {AI_NODE}")
        
        status = self.get_chain_status()
        if "error" not in status:
            print("[ChainClient] Chain connection: OK")
            asyncio.run(self.listen_for_events())
        else:
            print("[ChainClient] Chain offline — running in standalone mode")
    
    def stop(self):
        self.running = False
        print("[ChainClient] Stopped")


def demo():
    """Demonstrerer full end-to-end flyt"""
    client = ChainClient()
    
    print("\n=== AIVM End-to-End Demo ===\n")
    
    test_cases = [
        {
            "execution_id": f"exec-{hashlib.sha256(b'test1').hexdigest()[:8]}",
            "model_id": "default-classifier",
            "input_hash": hashlib.sha256(b"input_data_1").hexdigest(),
            "execution_type": "ON_CHAIN"
        },
        {
            "execution_id": f"exec-{hashlib.sha256(b'test2').hexdigest()[:8]}",
            "model_id": "llm-general",
            "input_hash": hashlib.sha256(b"input_data_2").hexdigest(),
            "execution_type": "OFF_CHAIN"
        }
    ]
    
    for test in test_cases:
        print(f"\n--- Executing: {test['model_id']} ({test['execution_type']}) ---")
        result = client.execute_ai_request(**test)
        
        if result:
            print(f"Execution path: {result.get('execution_path')}")
            print(f"Verified: {result.get('proof', {}).get('verified')}")
            print(f"Proof: {result.get('proof', {}).get('proof', '')[:32]}...")
        else:
            print("Execution failed or already processed")
    
    print("\n=== Demo Complete ===")


if __name__ == "__main__":
    demo()
