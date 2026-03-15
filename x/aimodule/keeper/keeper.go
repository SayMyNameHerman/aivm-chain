package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	logger   log.Logger
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, logger log.Logger) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		logger:   logger,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", types.ModuleName)
}

// --- AIModel CRUD ---

func (k Keeper) SetAIModel(ctx sdk.Context, model types.AIModel) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("model:%s", model.ModelId))
	bz, _ := json.Marshal(model)
	store.Set(key, bz)
}

func (k Keeper) GetAIModel(ctx sdk.Context, modelId string) (types.AIModel, bool) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("model:%s", modelId))
	bz := store.Get(key)
	if bz == nil {
		return types.AIModel{}, false
	}
	var model types.AIModel
	json.Unmarshal(bz, &model)
	return model, true
}

func (k Keeper) GetAllAIModels(ctx sdk.Context) []types.AIModel {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte("model:"))
	defer iterator.Close()
	var models []types.AIModel
	for ; iterator.Valid(); iterator.Next() {
		var model types.AIModel
		json.Unmarshal(iterator.Value(), &model)
		models = append(models, model)
	}
	return models
}

// --- ExecutionRequest CRUD ---

func (k Keeper) SetExecutionRequest(ctx sdk.Context, req types.ExecutionRequest) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("exec:%s", req.ExecutionId))
	bz, _ := json.Marshal(req)
	store.Set(key, bz)
}

func (k Keeper) GetExecutionRequest(ctx sdk.Context, execId string) (types.ExecutionRequest, bool) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("exec:%s", execId))
	bz := store.Get(key)
	if bz == nil {
		return types.ExecutionRequest{}, false
	}
	var req types.ExecutionRequest
	json.Unmarshal(bz, &req)
	return req, true
}
