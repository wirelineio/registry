//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emicklei/dot"
	"github.com/golang-collections/collections/stack"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Endpoints supported by the Querier.
const (
	ListResources = "list"
	GetResource   = "get"
	GetGraph      = "graph"
	GetTest       = "test"
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
		case GetTest:
			return getTest(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown utxo query endpoint.")
		}
	}
}

// nolint: unparam
func listResources(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records := keeper.ListResources(ctx)

	bz, err2 := json.MarshalIndent(records, "", "  ")
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

	bz, err2 := json.MarshalIndent(record, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getGraph(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	g := dot.NewGraph(dot.Directed)
	g.Attr("rankdir", "LR")

	if len(path) == 0 {
		resources := keeper.ListResources(ctx)
		for _, r := range resources {
			GraphResourceNode(g, r)
		}
	} else {
		pending := stack.New()
		done := make(map[string]bool)
		pending.Push(path[0])

		for pending.Len() > 0 {
			id := pending.Pop().(string)

			if _, exists := done[id]; !exists {
				r := keeper.GetResource(ctx, ID(id))
				GraphResourceNode(g, r)

				for link := range r.Links {
					pending.Push(link)
				}
			}

			done[id] = true
		}

	}

	return []byte(g.String()), nil
}

// nolint: unparam
func getTest(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	return []byte("test"), nil
}
