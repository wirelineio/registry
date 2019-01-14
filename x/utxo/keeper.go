//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	accountKeeper   auth.AccountKeeper
	coinKeeper      bank.Keeper
	accUtxoStoreKey sdk.StoreKey // Unexposed key to access Account UTXO store from sdk.Context.
	utxoStoreKey    sdk.StoreKey // Unexposed key to access UTXO store from sdk.Context.
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the UTXO Keeper.
func NewKeeper(accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, accUtxoStoreKey sdk.StoreKey, utxoStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		accountKeeper:   accountKeeper,
		coinKeeper:      coinKeeper,
		accUtxoStoreKey: accUtxoStoreKey,
		utxoStoreKey:    utxoStoreKey,
		cdc:             cdc,
	}
}

// PutAccOutput - saves an account UTXO to the store.
func (k Keeper) PutAccOutput(ctx sdk.Context, accUtxo AccOutput) {
	store := ctx.KVStore(k.accUtxoStoreKey)
	store.Set(accUtxo.ID, k.cdc.MustMarshalBinaryBare(accUtxo))
}

// ListAccOutput - get all account UTXO records.
func (k Keeper) ListAccOutput(ctx sdk.Context) []AccOutput {
	var records []AccOutput

	store := ctx.KVStore(k.accUtxoStoreKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj AccOutput
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			records = append(records, obj)
		}
	}

	return records
}

// PutOutPoint saves an outpoint to the UTXO store.
func (k Keeper) PutOutPoint(ctx sdk.Context, outpoint OutPoint) {
	store := ctx.KVStore(k.utxoStoreKey)
	store.Set([]byte(fmt.Sprintf("%s:%d", outpoint.Hash, outpoint.Index)), k.cdc.MustMarshalBinaryBare(outpoint))
}

// ListUtxo - get all account UTXO records.
func (k Keeper) ListUtxo(ctx sdk.Context) []OutPoint {
	var records []OutPoint

	store := ctx.KVStore(k.utxoStoreKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj OutPoint
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			records = append(records, obj)
		}
	}

	return records
}
