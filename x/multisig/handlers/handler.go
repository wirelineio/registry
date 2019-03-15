//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/registry/x/multisig/msgs"
)

// NewHandler returns a handler for "multisig" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case msgs.MsgInitMultiSig:
			return handleMsgInitMultiSig(ctx, keeper, msg)
		case msgs.MsgAbortMultiSig:
			return handleMsgAbortMultiSig(ctx, keeper, msg)
		case msgs.MsgJoinMultiSig:
			return handleMsgJoinMultiSig(ctx, keeper, msg)
		case msgs.MsgSpendMultiSig:
			return handleMsgSpendMultiSig(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized multisig Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
