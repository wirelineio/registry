//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"
	"strings"

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
	var namespace *string
	if len(path) > 0 {
		namespace = &path[0]
	}

	records := keeper.ListResources(ctx, namespace)

	bz, err2 := json.MarshalIndent(records, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getResource(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	id := ID(strings.Join(path, "/"))
	if !keeper.HasResource(ctx, id) {
		return nil, sdk.ErrInternal("Record not found.")
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
		records := keeper.ListResources(ctx, nil)
		for _, r := range records {
			GraphResourceNode(g, r)
		}
	} else {
		pending := stack.New()
		done := make(map[string]bool)
		id := strings.Join(path, "/")
		pending.Push(id)

		for pending.Len() > 0 {
			id := pending.Pop().(string)

			if _, exists := done[id]; !exists {
				r := keeper.GetResource(ctx, ID(id))
				GraphResourceNode(g, r)

				// for _, link := range r.Links {
				// 	if idAttr, ok := link["id"].(string); ok {
				// 		pending.Push(idAttr)
				// 	}

				// }
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
