import hashlib
import json
import time
from typing import Any, Dict
from sklearn.linear_model import LogisticRegression
from sklearn.preprocessing import StandardScaler
import numpy as np

class OnChainExecutor:
    """Kjører enkle AI-modeller direkte — transparent og deterministisk"""
    
    def __init__(self):
        self.models = {}
        self._load_default_models()
    
    def _load_default_models(self):
        # Enkel klassifiseringsmodell for demonstrasjon
        X = np.array([[1,2],[2,3],[3,4],[4,5],[5,6],[6,7]])
        y = np.array([0,0,0,1,1,1])
        
        scaler = StandardScaler()
        X_scaled = scaler.fit_transform(X)
        
        model = LogisticRegression()
        model.fit(X_scaled, y)
        
        self.models["default-classifier"] = {
            "model": model,
            "scaler": scaler,
            "type": "classification"
        }
    
    def execute(self, model_id: str, input_data: Dict) -> Dict:
        if model_id not in self.models and "classifier" in model_id:
            model_id = "default-classifier"
        
        if model_id not in self.models:
            return {"error": f"Model {model_id} not found", "status": "failed"}
        
        model_info = self.models[model_id]
        features = input_data.get("features", [1, 2])
        
        X = np.array([features])
        X_scaled = model_info["scaler"].transform(X)
        prediction = model_info["model"].predict(X_scaled)[0]
        probability = model_info["model"].predict_proba(X_scaled)[0].tolist()
        
        return {
            "model_id": model_id,
            "prediction": int(prediction),
            "probability": probability,
            "execution_type": "ON_CHAIN",
            "timestamp": int(time.time())
        }


class OffChainExecutor:
    """Kjører tyngre AI-modeller off-chain med proof-generering"""
    
    def __init__(self):
        self.execution_log = []
    
    def execute(self, model_id: str, input_data: Dict) -> Dict:
        # Simulerer tung AI-eksekvering
        time.sleep(0.1)
        
        input_str = json.dumps(input_data, sort_keys=True)
        result_hash = hashlib.sha256(
            f"{model_id}{input_str}{time.time()}".encode()
        ).hexdigest()
        
        result = {
            "model_id": model_id,
            "result_hash": result_hash,
            "execution_type": "OFF_CHAIN",
            "simulated_output": {
                "classification": "category_A",
                "confidence": 0.94,
                "processing_time_ms": 100
            },
            "timestamp": int(time.time())
        }
        
        self.execution_log.append({
            "model_id": model_id,
            "input_hash": hashlib.sha256(input_str.encode()).hexdigest(),
            "output_hash": result_hash,
            "timestamp": result["timestamp"]
        })
        
        return result
