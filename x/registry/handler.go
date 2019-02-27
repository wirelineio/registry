//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"
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

	payload := PayloadToPayloadYaml(msg.Payload)
	resource := payload.Resource

	exists := keeper.HasResource(ctx, resource.ID)

	if exists {
		// Check ownership.
		existingResource := keeper.GetResource(ctx, resource.ID)
		fmt.Println("Existing owner", existingResource.Owner.Address)
		fmt.Println("Signer", msg.Signer.String())
		if msg.Signer.String() != existingResource.Owner.Address {
			return sdk.ErrUnauthorized("Unauthorized resource write.").Result()
		}
	}

	keeper.PutResource(ctx, payload.Resource)

	bytes, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(bytes))

	return sdk.Result{}
}
