//
// Copyright 2018 Wireline, Inc.
//

package keepers

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	coinKeeper bank.Keeper

	multisigStoreKey sdk.StoreKey // Unexposed key to access HTLC store from sdk.Context.

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the multisig Keeper.
func NewKeeper(coinKeeper bank.Keeper, multisigStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:       coinKeeper,
		multisigStoreKey: multisigStoreKey,
		cdc:              cdc,
	}
}
