package cli

import (
	"github.com/spf13/cobra"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/keeper"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/types"
)

// getKeeper returns a transient keeper for CLI usage.
func getKeeper(cmd *cobra.Command) keeper.Keeper {
	clientCtx := client.GetClientContextFromCmd(cmd)
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	return keeper.NewKeeper(clientCtx.Codec, storeKey, log.NewNopLogger())
}
