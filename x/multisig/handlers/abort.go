//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wirechain/x/multisig/keepers"
	"github.com/wirelineio/wirechain/x/multisig/msgs"
)

// Handle MsgAbortMultiSig.
func handleMsgAbortMultiSig(ctx sdk.Context, keeper keepers.Keeper, msg msgs.MsgAbortMultiSig) sdk.Result {
	return sdk.Result{}
}
