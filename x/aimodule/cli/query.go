package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetQueryCmd returns the root aicli query command
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "Query AIVM AI module state",
	}

	cmd.AddCommand(
		CmdQueryModel(),
		CmdQueryAllModels(),
		CmdQueryExecution(),
		CmdQueryContract(),
		CmdQueryAllContracts(),
	)
	return cmd
}

// CmdQueryModel queries a specific AI model by ID
func CmdQueryModel() *cobra.Command {
	return &cobra.Command{
		Use:   "model [model-id]",
		Short: "Query a registered AI model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := sdk.Context{}
			k := getKeeper(cmd)
			model, found := k.GetAIModel(ctx, args[0])
			if !found {
				return fmt.Errorf("model '%s' not found", args[0])
			}
			out, _ := json.MarshalIndent(model, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
}

// CmdQueryAllModels queries all registered AI models
func CmdQueryAllModels() *cobra.Command {
	return &cobra.Command{
		Use:   "models",
		Short: "List all registered AI models",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := sdk.Context{}
			k := getKeeper(cmd)
			models := k.GetAllAIModels(ctx)
			if len(models) == 0 {
				fmt.Println("No models registered yet.")
				return nil
			}
			out, _ := json.MarshalIndent(models, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
}

// CmdQueryExecution queries a specific AI execution by ID
func CmdQueryExecution() *cobra.Command {
	return &cobra.Command{
		Use:   "execution [execution-id]",
		Short: "Query an AI execution request and its status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := sdk.Context{}
			k := getKeeper(cmd)
			req, found := k.GetExecutionRequest(ctx, args[0])
			if !found {
				return fmt.Errorf("execution '%s' not found", args[0])
			}
			out, _ := json.MarshalIndent(req, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
}

// CmdQueryContract queries a specific smart contract by ID
func CmdQueryContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [contract-id]",
		Short: "Query an uploaded Starlark smart contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := sdk.Context{}
			k := getKeeper(cmd)
			contract, found := k.GetSmartContract(ctx, args[0])
			if !found {
				return fmt.Errorf("contract '%s' not found", args[0])
			}
			showCode, _ := cmd.Flags().GetBool("show-code")
			if !showCode {
				contract.SourceCode = fmt.Sprintf("[%d chars] (use --show-code to display)", len(contract.SourceCode))
			}
			out, _ := json.MarshalIndent(contract, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().Bool("show-code", false, "Show the full Starlark source code")
	return cmd
}

// CmdQueryAllContracts lists all uploaded smart contracts
func CmdQueryAllContracts() *cobra.Command {
	return &cobra.Command{
		Use:   "contracts",
		Short: "List all uploaded Starlark smart contracts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := sdk.Context{}
			k := getKeeper(cmd)
			contracts := k.GetAllSmartContracts(ctx)
			if len(contracts) == 0 {
				fmt.Println("No contracts uploaded yet.")
				return nil
			}
			// Suppress source code for readability
			type summary struct {
				ContractId  string `json:"contract_id"`
				Owner       string `json:"owner"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Status      string `json:"status"`
				CreatedAt   int64  `json:"created_at"`
			}
			var summaries []summary
			for _, c := range contracts {
				summaries = append(summaries, summary{
					ContractId:  c.ContractId,
					Owner:       c.Owner,
					Name:        c.Name,
					Description: c.Description,
					Status:      c.Status,
					CreatedAt:   c.CreatedAt,
				})
			}
			out, _ := json.MarshalIndent(summaries, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
}
