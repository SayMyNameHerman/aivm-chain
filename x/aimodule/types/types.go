package types

const ModuleName = "aimodule"
const StoreKey = ModuleName

// ExecutionType definerer dual-path routing
const (
	ExecutionTypeOnChain  = "ON_CHAIN"
	ExecutionTypeOffChain = "OFF_CHAIN"
)

// ModelStatus
const (
	ModelStatusActive   = "ACTIVE"
	ModelStatusInactive = "INACTIVE"
)

// AIModel representerer en registrert AI-modell on-chain
type AIModel struct {
	ModelId       string `json:"model_id"`
	Owner         string `json:"owner"`
	ModelHash     string `json:"model_hash"`
	ExecutionType string `json:"execution_type"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"created_at"`
}

// ExecutionRequest representerer en AI-eksekveringsforespørsel
type ExecutionRequest struct {
	ExecutionId string `json:"execution_id"`
	ModelId     string `json:"model_id"`
	Requester   string `json:"requester"`
	InputHash   string `json:"input_hash"`
	OutputHash  string `json:"output_hash"`
	Proof       string `json:"proof"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}
