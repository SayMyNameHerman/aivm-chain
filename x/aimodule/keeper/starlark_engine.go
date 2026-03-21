package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/types"
	"go.starlark.net/starlark"
)

// =============================================================================
// Starlark Smart Contract Engine
// =============================================================================
// Contracts are short Python-like (Starlark) scripts stored on-chain.
// When executed, they run in a deterministic sandbox with access to two
// built-in functions:
//   - request_ai(model_id, input_data) → execution_id
//   - get_execution(execution_id)      → { status, output_hash, proof }
//
// Example contract:
//   def main(args):
//       exec_id = request_ai("my-model", args["prompt"])
//       return "Execution queued: " + exec_id
// =============================================================================

// StarlarkResult holds the output from a script execution
type StarlarkResult struct {
	ReturnValue string
	AIExecIds   []string
}

// ExecuteStarlark runs a Starlark script with AI built-ins injected into the environment.
func (k Keeper) ExecuteStarlark(
	ctx sdk.Context,
	contractId string,
	sourceCode string,
	caller string,
	args map[string]string,
) (StarlarkResult, error) {
	result := StarlarkResult{}
	var triggeredExecIds []string

	// --- Built-in: request_ai(model_id, input_data) ---
	requestAI := starlark.NewBuiltin("request_ai", func(
		thread *starlark.Thread,
		fn *starlark.Builtin,
		posArgs starlark.Tuple,
		kwargs []starlark.Tuple,
	) (starlark.Value, error) {
		if len(posArgs) < 2 {
			return nil, fmt.Errorf("request_ai requires 2 args: model_id, input_data")
		}
		modelID, ok := starlark.AsString(posArgs[0])
		if !ok {
			return nil, fmt.Errorf("request_ai: model_id must be a string")
		}
		inputData, ok := starlark.AsString(posArgs[1])
		if !ok {
			return nil, fmt.Errorf("request_ai: input_data must be a string")
		}

		// Hash of input_data acts as the inputHash for the on-chain record
		inputHash := fmt.Sprintf("%x", sha256.Sum256([]byte(inputData)))
		msgSrv := MsgServer{Keeper: k}
		execId, err := msgSrv.RequestExecution(ctx, caller, modelID, inputHash)
		if err != nil {
			return nil, fmt.Errorf("request_ai: %w", err)
		}

		triggeredExecIds = append(triggeredExecIds, execId)
		return starlark.String(execId), nil
	})

	// --- Built-in: get_execution(execution_id) ---
	getExecution := starlark.NewBuiltin("get_execution", func(
		thread *starlark.Thread,
		fn *starlark.Builtin,
		posArgs starlark.Tuple,
		kwargs []starlark.Tuple,
	) (starlark.Value, error) {
		if len(posArgs) < 1 {
			return nil, fmt.Errorf("get_execution requires 1 arg: execution_id")
		}
		execId, ok := starlark.AsString(posArgs[0])
		if !ok {
			return nil, fmt.Errorf("get_execution: execution_id must be a string")
		}

		req, found := k.GetExecutionRequest(ctx, execId)
		if !found {
			return starlark.None, nil
		}

		// Build a proper Starlark dict (implements starlark.Value)
		d := starlark.NewDict(4)
		d.SetKey(starlark.String("status"), starlark.String(req.Status))
		d.SetKey(starlark.String("output_hash"), starlark.String(req.OutputHash))
		d.SetKey(starlark.String("proof"), starlark.String(req.Proof))
		d.SetKey(starlark.String("model_id"), starlark.String(req.ModelId))
		return d, nil
	})

	// --- Build thread and predeclared environment ---
	thread := &starlark.Thread{Name: fmt.Sprintf("contract:%s", contractId)}
	predeclared := starlark.StringDict{
		"request_ai":    requestAI,
		"get_execution": getExecution,
	}

	// --- Execute the script ---
	globals, err := starlark.ExecFile(thread, contractId+".star", []byte(sourceCode), predeclared)
	if err != nil {
		return result, fmt.Errorf("starlark compile error: %w", err)
	}

	// --- Call the contract's main(args) function ---
	mainFn, ok := globals["main"]
	if !ok {
		return result, fmt.Errorf("contract must define a 'main(args)' function")
	}

	// Convert args to Starlark dict
	starlarkArgs := starlark.NewDict(len(args))
	for k, v := range args {
		starlarkArgs.SetKey(starlark.String(k), starlark.String(v))
	}

	returnVal, err := starlark.Call(thread, mainFn, starlark.Tuple{starlarkArgs}, nil)
	if err != nil {
		return result, fmt.Errorf("starlark runtime error: %w", err)
	}

	result.ReturnValue = returnVal.String()
	result.AIExecIds = triggeredExecIds
	return result, nil
}

// StoreContract stores a new smart contract on-chain
func (k Keeper) StoreContract(ctx context.Context, contract types.SmartContract) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	existing, found := k.GetSmartContract(sdkCtx, contract.ContractId)
	if found && existing.Status == types.ContractStatusActive {
		return fmt.Errorf("contract %s already exists", contract.ContractId)
	}
	contract.CreatedAt = time.Now().Unix()
	contract.Status = types.ContractStatusActive
	k.SetSmartContract(sdkCtx, contract)
	k.Logger().Info("Smart contract stored", "contract_id", contract.ContractId, "owner", contract.Owner)
	return nil
}

// RunContract retrieves and executes a smart contract by ID
func (k Keeper) RunContract(
	ctx sdk.Context,
	contractId string,
	caller string,
	args map[string]string,
) (types.ContractExecution, error) {
	contract, found := k.GetSmartContract(ctx, contractId)
	if !found {
		return types.ContractExecution{}, fmt.Errorf("contract %s not found", contractId)
	}
	if contract.Status != types.ContractStatusActive {
		return types.ContractExecution{}, fmt.Errorf("contract %s is not active", contractId)
	}

	starlarkResult, err := k.ExecuteStarlark(ctx, contractId, contract.SourceCode, caller, args)
	if err != nil {
		return types.ContractExecution{}, err
	}

	execution := types.ContractExecution{
		ContractId: contractId,
		Caller:     caller,
		Args:       args,
		Result:     starlarkResult.ReturnValue,
		AIExecIds:  starlarkResult.AIExecIds,
		ExecutedAt: time.Now().Unix(),
	}

	k.Logger().Info("Contract executed",
		"contract_id", contractId,
		"caller", caller,
		"result", starlarkResult.ReturnValue,
		"ai_executions", len(starlarkResult.AIExecIds),
	)

	return execution, nil
}
