package types

const ModuleName = "aimodule"
const StoreKey = ModuleName

// ExecutionType defines dual-path routing
const (
	ExecutionTypeOnChain  = "ON_CHAIN"
	ExecutionTypeOffChain = "OFF_CHAIN"
)

// ModelStatus
const (
	ModelStatusActive   = "ACTIVE"
	ModelStatusInactive = "INACTIVE"
)

// SmartContractStatus
const (
	ContractStatusActive   = "ACTIVE"
	ContractStatusDisabled = "DISABLED"
)

// AIModel represents a registered AI model on-chain
type AIModel struct {
	ModelId       string `json:"model_id"`
	Owner         string `json:"owner"`
	ModelHash     string `json:"model_hash"`
	ExecutionType string `json:"execution_type"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"created_at"`
}

// ExecutionRequest represents an AI execution request
type ExecutionRequest struct {
	ExecutionId string `json:"execution_id"`
	ModelId     string `json:"model_id"`
	Requester   string `json:"requester"`
	InputHash   string `json:"input_hash"`
	OutputHash  string `json:"output_hash"`
	Proof       string `json:"proof"`
	Status      string `json:"status"`
	ContractId  string `json:"contract_id,omitempty"` // Set if triggered by a smart contract
	CreatedAt   int64  `json:"created_at"`
}

// SmartContract represents an on-chain Starlark (Python) smart contract
type SmartContract struct {
	ContractId   string `json:"contract_id"`
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	// The Starlark (Python-like) source code of the contract
	SourceCode   string `json:"source_code"`
	Status       string `json:"status"`
	CreatedAt    int64  `json:"created_at"`
}

// ContractExecution represents the result of a smart contract execution
type ContractExecution struct {
	ContractId   string            `json:"contract_id"`
	Caller       string            `json:"caller"`
	Args         map[string]string `json:"args"`
	Result       string            `json:"result"`
	AIExecIds    []string          `json:"ai_exec_ids"` // All AI executions triggered by this contract
	ExecutedAt   int64             `json:"executed_at"`
}
