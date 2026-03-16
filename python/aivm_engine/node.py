import asyncio
import hashlib
import json
import time
import uuid
import requests
from typing import Dict, Optional
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from aivm_engine.executor import OnChainExecutor, OffChainExecutor
from aivm_engine.prover import create_execution_proof
from aivm_engine.router import ExecutionPath, get_routing_decision

app = FastAPI(title="AIVM Execution Node", version="1.0.0")

on_chain_executor = OnChainExecutor()
off_chain_executor = OffChainExecutor()

CHAIN_RPC = "http://localhost:26657"
CHAIN_REST = "http://localhost:1317"

class ExecutionRequest(BaseModel):
    execution_id: str
    model_id: str
    input_hash: str
    execution_type: str
    input_data: Dict

class ModelRegistration(BaseModel):
    model_id: str
    execution_type: str
    description: str = ""

@app.get("/health")
def health():
    return {
        "status": "online",
        "node": "AIVM Python Execution Engine",
        "version": "1.0.0",
        "timestamp": int(time.time())
    }

@app.post("/execute")
async def execute(request: ExecutionRequest):
    """
    Hovedendepunkt — mottar AI-eksekveringsforespørsler
    og returnerer kryptografisk bevis
    """
    path = get_routing_decision(
        request.model_id,
        request.execution_type,
        request.input_data
    )
    
    if path == ExecutionPath.ON_CHAIN:
        result = on_chain_executor.execute(request.model_id, request.input_data)
    else:
        result = off_chain_executor.execute(request.model_id, request.input_data)
    
    proof = create_execution_proof(
        execution_id=request.execution_id,
        model_id=request.model_id,
        input_hash=request.input_hash,
        output=result
    )
    
    return {
        "execution_id": request.execution_id,
        "execution_path": path.value,
        "result": result,
        "proof": {
            "output_hash": proof.output_hash,
            "proof": proof.proof,
            "timestamp": proof.timestamp,
            "verified": proof.verified
        }
    }

@app.get("/status/{execution_id}")
def get_status(execution_id: str):
    return {
        "execution_id": execution_id,
        "status": "completed",
        "timestamp": int(time.time())
    }

@app.get("/models")
def list_models():
    return {
        "on_chain_models": list(on_chain_executor.models.keys()),
        "off_chain_models": ["llm-general", "vision-classifier", "transformer-embedder"],
        "total": len(on_chain_executor.models) + 3
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
