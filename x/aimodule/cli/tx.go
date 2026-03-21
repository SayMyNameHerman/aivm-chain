package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmosregistry/chain-minimal/x/aimodule/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the root aicli tx command
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AIVM AI module transactions",
	}

	cmd.AddCommand(
		CmdRegisterModel(),
		CmdRequestExecution(),
		CmdSubmitProof(),
		CmdUploadContract(),
		CmdExecuteContract(),
	)
	return cmd
}

// CmdRegisterModel creates a CLI command to register an AI model on-chain
func CmdRegisterModel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-model [model-id] [model-hash] [execution-type] [creator-address]",
		Short: "Register an AI model on-chain (execution-type: ON_CHAIN or OFF_CHAIN)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			modelId := args[0]
			modelHash := args[1]
			execType := args[2]
			creator := args[3]

			ctx := sdk.Context{}
			k := getKeeper(cmd)
			msgSrv := keeper.NewMsgServer(k)
			id, err := msgSrv.RegisterModel(ctx, creator, modelId, modelHash, execType)
			if err != nil {
				return fmt.Errorf("failed to register model: %w", err)
			}
			fmt.Printf("✅ Model registered: %s\n", id)
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRequestExecution creates a CLI command to request an AI execution
func CmdRequestExecution() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-execution [model-id] [input-hash] [requester-address]",
		Short: "Request execution of a registered AI model",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			modelId := args[0]
			inputHash := args[1]
			requester := args[2]

			ctx := sdk.Context{}
			k := getKeeper(cmd)
			msgSrv := keeper.NewMsgServer(k)
			execId, err := msgSrv.RequestExecution(ctx, requester, modelId, inputHash)
			if err != nil {
				return fmt.Errorf("failed to request execution: %w", err)
			}
			fmt.Printf("✅ Execution requested. ID: %s\n", execId)
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitProof creates a CLI command to submit an AI execution proof
func CmdSubmitProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proof [execution-id] [output-hash] [proof]",
		Short: "Submit cryptographic proof for an off-chain AI execution",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			execId := args[0]
			outputHash := args[1]
			proof := args[2]

			ctx := sdk.Context{}
			k := getKeeper(cmd)
			msgSrv := keeper.NewMsgServer(k)
			success, err := msgSrv.SubmitProof(ctx, execId, outputHash, proof)
			if err != nil {
				return fmt.Errorf("failed to submit proof: %w", err)
			}
			if success {
				fmt.Printf("✅ Proof submitted and verified for: %s\n", execId)
			} else {
				fmt.Printf("❌ Proof verification failed\n")
			}
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdUploadContract creates a CLI command to upload a Starlark smart contract
func CmdUploadContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-contract [name] [description] [source-file-or-code]",
		Short: "Upload a Starlark (Python) smart contract on-chain",
		Long: `Upload a smart contract written in Starlark (Python dialect).

The contract must define a main(args) function. It has access to:
  - request_ai(model_id, input_data)  → returns execution_id
  - get_execution(execution_id)       → returns {status, output_hash, proof}

Example contract (save as contract.star):
  def main(args):
      exec_id = request_ai("my-model", args["prompt"])
      return "Queued: " + exec_id

Usage:
  aivmd ai upload-contract "My AI Contract" "Classifies prompts" contract.star`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			description := args[1]
			sourceArg := args[2]

			// If argument is a file path, read it; otherwise treat as inline code
			var sourceCode string
			if _, err := os.Stat(sourceArg); err == nil {
				bz, err := os.ReadFile(sourceArg)
				if err != nil {
					return fmt.Errorf("failed to read source file: %w", err)
				}
				sourceCode = string(bz)
			} else {
				sourceCode = sourceArg
			}

			owner, _ := cmd.Flags().GetString("from")
			contractId, _ := cmd.Flags().GetString("contract-id")

			ctx := sdk.Context{}
			k := getKeeper(cmd)
			msgSrv := keeper.NewMsgServer(k)
			id, err := msgSrv.UploadContract(ctx, owner, contractId, name, description, sourceCode)
			if err != nil {
				return fmt.Errorf("upload failed: %w", err)
			}
			fmt.Printf("✅ Contract uploaded!\n")
			fmt.Printf("   Contract ID: %s\n", id)
			fmt.Printf("   Name:        %s\n", name)
			fmt.Printf("   Lines:       %d\n", len(strings.Split(sourceCode, "\n")))
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("contract-id", "", "Custom contract ID (auto-generated if omitted)")
	return cmd
}

// CmdExecuteContract creates a CLI command to execute a Starlark smart contract
func CmdExecuteContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-contract [contract-id]",
		Short: "Execute an on-chain Starlark smart contract",
		Long: `Execute a previously uploaded Starlark smart contract.

Pass arguments as key=value pairs using --arg flags.

Example:
  aivmd ai execute-contract my-contract-id --from alice --arg prompt="Hello AI"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contractId := args[0]
			caller, _ := cmd.Flags().GetString("from")
			argFlags, _ := cmd.Flags().GetStringArray("arg")

			// Parse key=value args
			contractArgs := make(map[string]string)
			for _, a := range argFlags {
				parts := strings.SplitN(a, "=", 2)
				if len(parts) == 2 {
					contractArgs[parts[0]] = parts[1]
				}
			}

			ctx := sdk.Context{}
			k := getKeeper(cmd)
			msgSrv := keeper.NewMsgServer(k)
			execution, err := msgSrv.ExecuteContract(ctx, caller, contractId, contractArgs)
			if err != nil {
				return fmt.Errorf("execution failed: %w", err)
			}

			out, _ := json.MarshalIndent(execution, "", "  ")
			fmt.Printf("✅ Contract executed!\n\n%s\n", string(out))
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().StringArray("arg", []string{}, "Contract argument (key=value, can repeat)")
	return cmd
}
