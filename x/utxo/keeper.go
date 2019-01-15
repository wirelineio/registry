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
	txStoreKey      sdk.StoreKey // Unexposed key to access TX store from sdk.Context.
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the UTXO Keeper.
func NewKeeper(accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, accUtxoStoreKey sdk.StoreKey, utxoStoreKey sdk.StoreKey, txStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		accountKeeper:   accountKeeper,
		coinKeeper:      coinKeeper,
		accUtxoStoreKey: accUtxoStoreKey,
		utxoStoreKey:    utxoStoreKey,
		txStoreKey:      txStoreKey,
		cdc:             cdc,
	}
}

// PutAccOutput - saves an account UTXO to the store.
func (k Keeper) PutAccOutput(ctx sdk.Context, accUtxo AccOutput) {
	store := ctx.KVStore(k.accUtxoStoreKey)
	store.Set(accUtxo.ID, k.cdc.MustMarshalBinaryBare(accUtxo))
}

// HasAccOutput - checks if an account output by the given ID exists.
func (k Keeper) HasAccOutput(ctx sdk.Context, id []byte) bool {
	store := ctx.KVStore(k.accUtxoStoreKey)
	return store.Has(id)
}

// GetAccOutput - gets a AccOutput from the store.
func (k Keeper) GetAccOutput(ctx sdk.Context, id []byte) AccOutput {
	store := ctx.KVStore(k.accUtxoStoreKey)

	bz := store.Get(id)
	var obj AccOutput
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
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

// GetOutPointKey returns the key used in the KVStore for the given OutPoint.
func GetOutPointKey(op OutPoint) string {
	return fmt.Sprintf("%s:%d", op.Hash, op.Index)
}

// PutOutPoint saves an outpoint to the UTXO store.
func (k Keeper) PutOutPoint(ctx sdk.Context, outpoint OutPoint) {
	store := ctx.KVStore(k.utxoStoreKey)
	store.Set([]byte(GetOutPointKey(outpoint)), k.cdc.MustMarshalBinaryBare(outpoint))
}

// HasOutPoint checks if the given outpoint exists in the UTXO list.
func (k Keeper) HasOutPoint(ctx sdk.Context, outpoint OutPoint) bool {
	store := ctx.KVStore(k.utxoStoreKey)
	return store.Has([]byte(GetOutPointKey(outpoint)))
}

// DeleteOutPoint deletes the given outpoint from the UTXO list.
func (k Keeper) DeleteOutPoint(ctx sdk.Context, outpoint OutPoint) {
	store := ctx.KVStore(k.utxoStoreKey)
	store.Delete([]byte(GetOutPointKey(outpoint)))
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

// HasTx checks if a transaction by the given hash exists.
func (k Keeper) HasTx(ctx sdk.Context, hash []byte) bool {
	store := ctx.KVStore(k.txStoreKey)
	return store.Has(hash)
}

// PutTx - saves a transaction to the store.
func (k Keeper) PutTx(ctx sdk.Context, hash []byte, tx Tx) {
	store := ctx.KVStore(k.txStoreKey)
	store.Set(hash, k.cdc.MustMarshalBinaryBare(tx))
}

// GetTx - gets a transaction from the store.
func (k Keeper) GetTx(ctx sdk.Context, hash []byte) Tx {
	store := ctx.KVStore(k.txStoreKey)

	bz := store.Get(hash)
	var obj Tx
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
}
