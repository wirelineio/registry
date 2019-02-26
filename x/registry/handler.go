//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "utxo" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetResource:
			return handleMsgSetResource(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized registry Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetResource.
func handleMsgSetResource(ctx sdk.Context, keeper Keeper, msg MsgSetResource) sdk.Result {
	fmt.Println("---------------------------- handleMsgSetResource -----------------------------")

	return sdk.Result{}
}
