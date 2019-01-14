//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Endpoints supported by the Querier.
const (
	ListAccOutput = "ls-account-outputs"
	ListUtxo      = "ls"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListAccOutput:
			return listAccOutput(ctx, path[1:], req, keeper)
		case ListUtxo:
			return listUtxo(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown utxo query endpoint.")
		}
	}
}

// nolint: unparam
func listAccOutput(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records := keeper.ListAccOutput(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, records)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func listUtxo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records := keeper.ListUtxo(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, records)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}
