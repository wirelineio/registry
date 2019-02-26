//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emicklei/dot"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Endpoints supported by the Querier.
const (
	ListResources = "list"
	GetResource   = "get"
	GetGraph      = "graph"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListResources:
			return listResources(ctx, path[1:], req, keeper)
		case GetResource:
			return getResource(ctx, path[1:], req, keeper)
		case GetGraph:
			return getGraph(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown utxo query endpoint.")
		}
	}
}

// nolint: unparam
func listResources(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records := keeper.ListResources(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, records)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getResource(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	id := ID(path[0])
	if !keeper.HasResource(ctx, id) {
		return nil, sdk.ErrInternal("Resource not found.")
	}

	record := keeper.GetResource(ctx, id)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, record)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getGraph(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	g := dot.NewGraph(dot.Directed)

	return []byte(g.String()), nil
}
