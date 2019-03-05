//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

// NewHandler returns a handler for "utxo" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetResource:
			return handleMsgSetResource(ctx, keeper, msg)
		case MsgDeleteResource:
			return handleMsgDeleteResource(ctx, keeper, msg)
		case MsgClearResources:
			return handleMsgClearResources(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized registry Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetResource.
func handleMsgSetResource(ctx sdk.Context, keeper Keeper, msg MsgSetResource) sdk.Result {
	payload := PayloadObjToPayload(msg.Payload)
	resource := payload.Resource

	if exists := keeper.HasResource(ctx, resource.ID); exists {
		// Check ownership.
		existingResource := keeper.GetResource(ctx, resource.ID)

		allow := checkAccess(existingResource, payload.Signatures)
		if !allow {
			return sdk.ErrUnauthorized("Unauthorized resource write.").Result()
		}
	}

	keeper.PutResource(ctx, payload.Resource)

	return sdk.Result{}
}

// Handle MsgDeleteResource.
func handleMsgDeleteResource(ctx sdk.Context, keeper Keeper, msg MsgDeleteResource) sdk.Result {
	payload := PayloadObjToPayload(msg.Payload)
	resource := payload.Resource

	if exists := keeper.HasResource(ctx, resource.ID); exists {
		// Check ownership.
		existingResource := keeper.GetResource(ctx, resource.ID)

		allow := checkAccess(existingResource, payload.Signatures)
		if !allow {
			return sdk.ErrUnauthorized("Unauthorized resource write.").Result()
		}

		keeper.DeleteResource(ctx, payload.Resource.ID)

		return sdk.Result{}
	}

	return sdk.ErrInternal("Resource not found.").Result()
}

// Handle MsgClearResources.
func handleMsgClearResources(ctx sdk.Context, keeper Keeper, msg MsgClearResources) sdk.Result {
	keeper.ClearResources(ctx)

	return sdk.Result{}
}

func checkAccess(resource Resource, signatures []Signature) bool {
	addresses := make(map[string]bool)

	// Check signatures.
	resourceSignBytes := GenResourceHash(resource)
	for _, sig := range signatures {
		pubKey, err := cryptoAmino.PubKeyFromBytes(BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes.")
			return false
		}

		addresses[GetAddressFromPubKey(pubKey)] = true

		allow := pubKey.VerifyBytes(resourceSignBytes, BytesFromBase64(sig.Signature))
		if !allow {
			fmt.Println("Signature mismatch: ", sig.PubKey)

			return false
		}
	}

	// Check one of the addresses matches the owner.
	_, ok := addresses[resource.Owner.Address]
	if !ok {
		return false
	}

	return true
}
