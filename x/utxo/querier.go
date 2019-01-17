//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"encoding/base64"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/emicklei/dot"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wirechain/x/utxo/utils"
)

// Endpoints supported by the Querier.
const (
	ListAccOutput = "ls-account-outputs"
	ListUtxo      = "ls"
	ListTx        = "ls-tx"
	GetTx         = "get-tx"
	GetBalance    = "balance"
	GetGraph      = "graph"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListAccOutput:
			return listAccOutput(ctx, path[1:], req, keeper)
		case ListUtxo:
			return listUtxo(ctx, path[1:], req, keeper)
		case ListTx:
			return listTx(ctx, path[1:], req, keeper)
		case GetTx:
			return getTx(ctx, path[1:], req, keeper)
		case GetBalance:
			return getBalance(ctx, path[1:], req, keeper)
		case GetGraph:
			return getGraph(ctx, path[1:], req, keeper)
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

// nolint: unparam
func listTx(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records, _ := keeper.ListTx(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, records)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getTx(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	// strings.Join works around an issue where the base64 encoded string contains a slash.
	hash, err2 := base64.StdEncoding.DecodeString(strings.Join(path, "/"))

	if err2 != nil {
		return nil, sdk.ErrInternal("Error decoding transaction hash.")
	}

	if !keeper.HasTx(ctx, hash) {
		return nil, sdk.ErrInternal("Transaction not found.")
	}

	record := keeper.GetTx(ctx, hash)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, record)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getBalance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdk.ErrInvalidAddress(path[0])
	}

	var wallet Wallet

	// For each UTXO:
	// Get the transaction output.
	// Check if it's payable to the given address.
	// If so, add UTXO value to current balance.
	utxos := keeper.ListUtxo(ctx)
	for _, outpoint := range utxos {
		if outpoint.Index >= 0 {
			tx := keeper.GetTx(ctx, outpoint.Hash)
			txOut := tx.TxOut[outpoint.Index]

			var obj PayToAddress
			keeper.cdc.MustUnmarshalBinaryBare(txOut.PkScript, &obj)
			if obj.Address.Equals(address) {
				updateWallet(&wallet, outpoint, txOut.Value)
			}

		} else if outpoint.Index == OutPointAccountBirth {
			accOutput := keeper.GetAccOutput(ctx, outpoint.Hash)
			if accOutput.Address.Equals(address) {
				updateWallet(&wallet, outpoint, accOutput.Value)
			}
		}
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, wallet)
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

func updateWallet(wallet *Wallet, outpoint OutPoint, value uint64) {
	wallet.Balance += value
	wallet.Entries = append(wallet.Entries, OutPointVal{
		Hash:  outpoint.Hash,
		Index: outpoint.Index,
		Value: value,
	})
}

// nolint: unparam
func getGraph(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	g := dot.NewGraph(dot.Directed)

	txns, txIds := keeper.ListTx(ctx)
	for txIndex, tx := range txns {
		txNodeID := utils.BytesToBase64(txIds[txIndex])
		var txLabel []string
		txDotNode := TxNode(g, txNodeID)

		txLabel = append(txLabel, "TXN")
		txLabel = append(txLabel, txNodeID)

		for _, txIn := range tx.TxIn {
			g.Edge(
				g.Node(txNodeID),
				TxNode(g, utils.BytesToBase64(txIn.Input.Hash)),
				TxInLabel(txIn),
			)
		}

		for txOutputIndex, txOut := range tx.TxOut {
			txLabel = append(txLabel, TxOutLabel(txOutputIndex, txOut))
		}

		txDotNode.Attr("label", strings.Join(txLabel, " | "))
	}

	for _, accOut := range keeper.ListAccOutput(ctx) {
		AccOutNode(g, accOut)
	}

	for _, utxo := range keeper.ListUtxo(ctx) {
		UnspentOutputNode(g, utxo)
	}

	return []byte(g.String()), nil
}
