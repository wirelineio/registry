package htlc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	coinKeeper bank.Keeper

	htlcStoreKey sdk.StoreKey // Unexposed key to access HTLC store from sdk.Context.

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper.
func NewKeeper(coinKeeper bank.Keeper, htlcStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:   coinKeeper,
		htlcStoreKey: htlcStoreKey,
		cdc:          cdc,
	}
}

// HasHtlc - returns whether or not the HTLC by that hash exists.
func (k Keeper) HasHtlc(ctx sdk.Context, hash string) bool {
	store := ctx.KVStore(k.htlcStoreKey)
	bz := store.Get([]byte(hash))
	return bz != nil
}
