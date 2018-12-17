//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// State of the contract.
type State int8

// Contract state enum.
const (
	StateCreated State = 1
	StateLocked  State = 2
	StateAborted State = 3
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	coinKeeper bank.Keeper

	multisigStoreKey sdk.StoreKey // Unexposed key to access HTLC store from sdk.Context.

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// Contract represents the state of the contract.
type Contract struct {
	ID           string
	State        State
	AliceAmount  sdk.Coin
	AliceAddress sdk.AccAddress
	BobAmount    sdk.Coin
	BobAddress   sdk.AccAddress
	Balance      sdk.Coin
}

// NewKeeper creates new instances of the multisig Keeper.
func NewKeeper(coinKeeper bank.Keeper, multisigStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:       coinKeeper,
		multisigStoreKey: multisigStoreKey,
		cdc:              cdc,
	}
}

// HasContract - returns whether or not a contract by that ID exists.
func (k Keeper) HasContract(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.multisigStoreKey)
	bz := store.Get([]byte(id))
	return bz != nil
}

// UpsertContract - inserts/updates contract.
func (k Keeper) UpsertContract(ctx sdk.Context, obj Contract) {
	store := ctx.KVStore(k.multisigStoreKey)
	store.Set([]byte(obj.ID), k.cdc.MustMarshalBinaryBare(obj))
}
