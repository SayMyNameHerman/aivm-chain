package module

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/keeper"
	"github.com/cosmosregistry/chain-minimal/x/aimodule/types"
)

const ConsensusVersion = 1

var (
	_ module.AppModuleBasic = AppModule{}
	_ appmodule.AppModule   = AppModule{}
)

type AppModule struct {
	keeper keeper.Keeper
}

func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (AppModule) Name() string                                         { return types.ModuleName }
func (AppModule) IsOnePerModuleType()                                  {}
func (AppModule) IsAppModule()                                         {}
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino)         {}
func (AppModule) RegisterInterfaces(codectypes.InterfaceRegistry)      {}
func (AppModule) RegisterGRPCGatewayRoutes(client.Context, *gwruntime.ServeMux) {}
func (AppModule) DefaultGenesis() json.RawMessage                      { return json.RawMessage(`{}`) }
func (AppModule) ValidateGenesis(json.RawMessage) error                { return nil }
func (am AppModule) InitGenesis(ctx context.Context, data json.RawMessage) error { return nil }
func (am AppModule) ExportGenesis(ctx context.Context) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }
func (am AppModule) BeginBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	am.keeper.Logger().Info("aimodule BeginBlock", "height", sdkCtx.BlockHeight())
	return nil
}
