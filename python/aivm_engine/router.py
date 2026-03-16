from typing import Dict, Tuple
from enum import Enum

class ExecutionPath(Enum):
    ON_CHAIN = "ON_CHAIN"
    OFF_CHAIN = "OFF_CHAIN"

# Modeller som kjøres on-chain (enkle, deterministiske)
ON_CHAIN_MODELS = {
    "classifier",
    "scoring",
    "risk-assessment",
    "governance",
}

# Modeller som rutes off-chain (tunge, komplekse)
OFF_CHAIN_MODELS = {
    "llm",
    "vision",
    "transformer",
    "embedding",
}

def route_execution(model_id: str, input_data: Dict) -> Tuple[ExecutionPath, str]:
    """
    Dual-path routing — kjernelogikken i AIVM-arkitekturen.
    Bestemmer om en AI-jobb skal kjøres on-chain eller off-chain.
    """
    model_id_lower = model_id.lower()
    
    # Sjekk eksplisitte on-chain modeller
    for keyword in ON_CHAIN_MODELS:
        if keyword in model_id_lower:
            return ExecutionPath.ON_CHAIN, f"Model '{model_id}' matches on-chain keyword '{keyword}'"
    
    # Sjekk eksplisitte off-chain modeller
    for keyword in OFF_CHAIN_MODELS:
        if keyword in model_id_lower:
            return ExecutionPath.OFF_CHAIN, f"Model '{model_id}' matches off-chain keyword '{keyword}'"
    
    # Rut basert på input-kompleksitet
    input_size = len(str(input_data))
    if input_size > 1000:
        return ExecutionPath.OFF_CHAIN, f"Input size {input_size} bytes exceeds on-chain threshold"
    
    # Default: on-chain for enkle modeller
    return ExecutionPath.ON_CHAIN, "Default routing: simple model runs on-chain"

def get_routing_decision(model_id: str, execution_type: str, input_data: Dict) -> ExecutionPath:
    """
    Respekterer chain-registrert executionType, men kan override
    basert på faktisk input-kompleksitet.
    """
    if execution_type == "OFF_CHAIN":
        return ExecutionPath.OFF_CHAIN
    
    if execution_type == "ON_CHAIN":
        # Kan override til off-chain hvis input er for stort
        input_size = len(str(input_data))
        if input_size > 5000:
            return ExecutionPath.OFF_CHAIN
        return ExecutionPath.ON_CHAIN
    
    # Ukjent type — bruk automatisk routing
    path, _ = route_execution(model_id, input_data)
    return path
