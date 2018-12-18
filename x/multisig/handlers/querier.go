//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Endpoints supported by the Querier.
const (
	QueryView = "view"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryView:
			return queryView(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown multisig query endpoint.")
		}
	}
}

// nolint: unparam
func queryView(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	id := path[0]

	if !keeper.HasContract(ctx, id) {
		return []byte{}, sdk.ErrUnknownRequest("Contract not found.")
	}

	value := keeper.GetContract(ctx, id)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, value)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}
