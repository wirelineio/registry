//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wirechain/x/multisig/msgs"
)

// Handle MsgAbortMultiSig.
func handleMsgAbortMultiSig(ctx sdk.Context, keeper Keeper, msg msgs.MsgAbortMultiSig) sdk.Result {
	if !keeper.HasContract(ctx, msg.ID) {
		return sdk.ErrInternal("Contract not found.").Result()
	}

	obj := keeper.GetContract(ctx, msg.ID)
	if !obj.AliceAddress.Equals(msg.AliceAddress) {
		return sdk.ErrInternal("Message signer does not own the contract.").Result()
	}

	if obj.State != StateCreated {
		return sdk.ErrInternal("Contract cannot be aborted in the present state.").Result()
	}

	// Refund the coins to Alice.
	_, _, err := keeper.coinKeeper.AddCoins(ctx, obj.AliceAddress, sdk.Coins{obj.AliceAmount})
	if err != nil {
		return sdk.ErrInsufficientCoins("Error returning coins.").Result()
	}

	keeper.DeleteContract(ctx, msg.ID)

	return sdk.Result{}
}
