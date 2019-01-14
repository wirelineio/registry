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
	ListAccountUtxo = "ls-account-utxo"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListAccountUtxo:
			return listAccountUtxo(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown utxo query endpoint.")
		}
	}
}

// nolint: unparam
func listAccountUtxo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records := keeper.ListAccountUtxo(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, records)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}
