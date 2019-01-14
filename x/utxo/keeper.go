//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
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

// PutAccountUtxo - saves an account UTXO to the store.
func (k Keeper) PutAccountUtxo(ctx sdk.Context, accUtxo AccountUtxo) {
	store := ctx.KVStore(k.accUtxoStoreKey)
	store.Set(accUtxo.ID, k.cdc.MustMarshalBinaryBare(accUtxo))
}

// ListAccountUtxo - get all account UTXO records.
func (k Keeper) ListAccountUtxo(ctx sdk.Context) []AccountUtxo {
	var records []AccountUtxo

	store := ctx.KVStore(k.accUtxoStoreKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj AccountUtxo
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			records = append(records, obj)
		}
	}

	return records
}
