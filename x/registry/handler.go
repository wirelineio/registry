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
		case MsgSetRecord:
			return handleMsgSetResource(ctx, keeper, msg)
		case MsgDeleteRecord:
			return handleMsgDeleteResource(ctx, keeper, msg)
		case MsgClearRecords:
			return handleMsgClearResources(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized registry Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetRecord.
func handleMsgSetResource(ctx sdk.Context, keeper Keeper, msg MsgSetRecord) sdk.Result {
	payload := PayloadObjToPayload(msg.Payload)
	record := payload.Record

	if exists := keeper.HasResource(ctx, record.ID); exists {
		// Check ownership.
		owner := keeper.GetResource(ctx, record.ID).Owner

		allow := checkAccess(owner, record, payload.Signatures)
		if !allow {
			return sdk.ErrUnauthorized("Unauthorized record write.").Result()
		}
	}

	keeper.PutResource(ctx, payload.Record)

	return sdk.Result{}
}

// Handle MsgDeleteRecord.
func handleMsgDeleteResource(ctx sdk.Context, keeper Keeper, msg MsgDeleteRecord) sdk.Result {
	payload := PayloadObjToPayload(msg.Payload)
	record := payload.Record

	if exists := keeper.HasResource(ctx, record.ID); exists {
		// Check ownership.
		owner := keeper.GetResource(ctx, record.ID).Owner

		allow := checkAccess(owner, record, payload.Signatures)
		if !allow {
			return sdk.ErrUnauthorized("Unauthorized record write.").Result()
		}

		keeper.DeleteResource(ctx, payload.Record.ID)

		return sdk.Result{}
	}

	return sdk.ErrInternal("Record not found.").Result()
}

// Handle MsgClearRecords.
func handleMsgClearResources(ctx sdk.Context, keeper Keeper, msg MsgClearRecords) sdk.Result {
	keeper.ClearResources(ctx)

	return sdk.Result{}
}

func checkAccess(owner string, record Record, signatures []Signature) bool {
	addresses := make(map[string]bool)

	// Check signatures.
	resourceSignBytes := GenRecordHash(record)
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
	_, ok := addresses[owner]
	if !ok {
		return false
	}

	return true
}
