package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/types"
)

type MsgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) MsgServer {
	return MsgServer{Keeper: keeper}
}

// RegisterModel registrerer en AI-modell on-chain
func (m MsgServer) RegisterModel(ctx context.Context, creator, modelId, modelHash, executionType string) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if executionType != types.ExecutionTypeOnChain && executionType != types.ExecutionTypeOffChain {
		return "", fmt.Errorf("ugyldig executionType: må være ON_CHAIN eller OFF_CHAIN")
	}

	_, found := m.GetAIModel(sdkCtx, modelId)
	if found {
		return "", fmt.Errorf("modell %s eksisterer allerede", modelId)
	}

	model := types.AIModel{
		ModelId:       modelId,
		Owner:         creator,
		ModelHash:     modelHash,
		ExecutionType: executionType,
		Status:        types.ModelStatusActive,
		CreatedAt:     time.Now().Unix(),
	}

	m.SetAIModel(sdkCtx, model)
	m.Logger().Info("Modell registrert", "modelId", modelId, "type", executionType)

	return modelId, nil
}

// RequestExecution — DUAL-PATH ROUTING KJERNELOGIKK
func (m MsgServer) RequestExecution(ctx context.Context, requester, modelId, inputHash string) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	model, found := m.GetAIModel(sdkCtx, modelId)
	if !found {
		return "", fmt.Errorf("modell %s ikke funnet", modelId)
	}

	execId := generateExecId(modelId, inputHash, sdkCtx.BlockHeight())

	var status string
	switch model.ExecutionType {
	case types.ExecutionTypeOnChain:
		// Enkel modell — kjør direkte on-chain
		status = "EXECUTED_ON_CHAIN"
		m.Logger().Info("ON_CHAIN eksekvering", "execId", execId, "modelId", modelId)

	case types.ExecutionTypeOffChain:
		// Kompleks modell — rut til off-chain node
		status = "PENDING_OFF_CHAIN"
		m.Logger().Info("OFF_CHAIN routing", "execId", execId, "modelId", modelId)
	}

	req := types.ExecutionRequest{
		ExecutionId: execId,
		ModelId:     modelId,
		Requester:   requester,
		InputHash:   inputHash,
		Status:      status,
		CreatedAt:   time.Now().Unix(),
	}

	m.SetExecutionRequest(sdkCtx, req)
	return execId, nil
}

// SubmitProof — off-chain node leverer kryptografisk bevis
func (m MsgServer) SubmitProof(ctx context.Context, execId, outputHash, proof string) (bool, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	req, found := m.GetExecutionRequest(sdkCtx, execId)
	if !found {
		return false, fmt.Errorf("execution %s ikke funnet", execId)
	}

	if req.Status != "PENDING_OFF_CHAIN" {
		return false, fmt.Errorf("execution %s er ikke i PENDING_OFF_CHAIN status", execId)
	}

	// Verifiser proof — hash av execId + outputHash
	expectedProof := fmt.Sprintf("%x", sha256.Sum256([]byte(execId+outputHash)))
	if proof != expectedProof {
		return false, fmt.Errorf("ugyldig proof for execution %s", execId)
	}

	req.OutputHash = outputHash
	req.Proof = proof
	req.Status = "VERIFIED"
	m.SetExecutionRequest(sdkCtx, req)

	m.Logger().Info("Proof verifisert", "execId", execId, "outputHash", outputHash)
	return true, nil
}

func generateExecId(modelId, inputHash string, height int64) string {
	data := fmt.Sprintf("%s-%s-%d-%d", modelId, inputHash, height, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:16])
}
