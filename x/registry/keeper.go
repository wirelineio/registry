//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	accountKeeper    auth.AccountKeeper
	coinKeeper       bank.Keeper
	resourceStoreKey sdk.StoreKey // Unexposed key to access record store from sdk.Context.
	cdc              *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the UTXO Keeper.
func NewKeeper(accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, resourceStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		accountKeeper:    accountKeeper,
		coinKeeper:       coinKeeper,
		resourceStoreKey: resourceStoreKey,
		cdc:              cdc,
	}
}

// PutResource - saves a record to the store.
func (k Keeper) PutResource(ctx sdk.Context, record Record) {
	store := ctx.KVStore(k.resourceStoreKey)
	store.Set([]byte(record.ID), k.cdc.MustMarshalBinaryBare(RecordToRecordObj(record)))
}

// HasResource - checks if a record by the given ID exists.
func (k Keeper) HasResource(ctx sdk.Context, id ID) bool {
	store := ctx.KVStore(k.resourceStoreKey)
	return store.Has([]byte(id))
}

// GetResource - gets a record from the store.
func (k Keeper) GetResource(ctx sdk.Context, id ID) Record {
	store := ctx.KVStore(k.resourceStoreKey)

	bz := store.Get([]byte(id))
	var obj RecordObj
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return RecordObjToRecord(obj)
}

// ListResources - get all record records.
func (k Keeper) ListResources(ctx sdk.Context, namespace *string) []Record {
	var records []Record

	store := ctx.KVStore(k.resourceStoreKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj RecordObj
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)

			record := RecordObjToRecord(obj)
			if namespace == nil {
				records = append(records, record)
			} else if ns, ok := record.Attributes["namespace"].(string); ok && *namespace == ns {
				records = append(records, record)
			}
		}
	}

	return records
}

// DeleteResource - deletes a record from the store.
func (k Keeper) DeleteResource(ctx sdk.Context, id ID) {
	store := ctx.KVStore(k.resourceStoreKey)
	store.Delete([]byte(id))
}

// ClearResources - Deletes all records.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func (k Keeper) ClearResources(ctx sdk.Context) {
	store := ctx.KVStore(k.resourceStoreKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		store.Delete(itr.Key())
	}
}
