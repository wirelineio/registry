//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wirechain/x/multisig/msgs"
)

// Handle MsgJoinMultiSig.
func handleMsgJoinMultiSig(ctx sdk.Context, keeper Keeper, msg msgs.MsgJoinMultiSig) sdk.Result {
	return sdk.Result{}
}
