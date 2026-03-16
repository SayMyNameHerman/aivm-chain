import hashlib
import hmac
import json
import time
from dataclasses import dataclass

@dataclass
class ExecutionProof:
    execution_id: str
    model_id: str
    input_hash: str
    output_hash: str
    proof: str
    timestamp: int
    verified: bool

def hash_input(data: str) -> str:
    return hashlib.sha256(data.encode()).hexdigest()

def hash_output(output: any) -> str:
    output_str = json.dumps(output, sort_keys=True)
    return hashlib.sha256(output_str.encode()).hexdigest()

def generate_proof(execution_id: str, output_hash: str) -> str:
    """Genererer kryptografisk bevis — samme logikk som Go-chain forventer"""
    data = f"{execution_id}{output_hash}"
    return hashlib.sha256(data.encode()).hexdigest()

def verify_proof(execution_id: str, output_hash: str, proof: str) -> bool:
    expected = generate_proof(execution_id, output_hash)
    return hmac.compare_digest(expected, proof)

def create_execution_proof(execution_id: str, model_id: str, 
                           input_hash: str, output: any) -> ExecutionProof:
    output_hash = hash_output(output)
    proof = generate_proof(execution_id, output_hash)
    
    return ExecutionProof(
        execution_id=execution_id,
        model_id=model_id,
        input_hash=input_hash,
        output_hash=output_hash,
        proof=proof,
        timestamp=int(time.time()),
        verified=True
    )
