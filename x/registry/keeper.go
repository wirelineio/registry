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
	resourceStoreKey sdk.StoreKey // Unexposed key to access resource store from sdk.Context.
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

// PutResource - saves a resource to the store.
func (k Keeper) PutResource(ctx sdk.Context, resource Resource) {
	store := ctx.KVStore(k.resourceStoreKey)
	store.Set([]byte(resource.ID), k.cdc.MustMarshalBinaryBare(resource))
}

// HasResource - checks if a resource by the given ID exists.
func (k Keeper) HasResource(ctx sdk.Context, id ID) bool {
	store := ctx.KVStore(k.resourceStoreKey)
	return store.Has([]byte(id))
}

// GetResource - gets a resource from the store.
func (k Keeper) GetResource(ctx sdk.Context, id ID) Resource {
	store := ctx.KVStore(k.resourceStoreKey)

	bz := store.Get([]byte(id))
	var obj Resource
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
}

// ListResources - get all resource records.
func (k Keeper) ListResources(ctx sdk.Context) []Resource {
	var records []Resource

	store := ctx.KVStore(k.resourceStoreKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj Resource
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			records = append(records, obj)
		}
	}

	return records
}
