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
	if !keeper.HasContract(ctx, msg.ID) {
		return sdk.ErrInternal("Contract not found.").Result()
	}

	obj := keeper.GetContract(ctx, msg.ID)
	if !obj.BobAddress.Equals(msg.BobAddress) {
		return sdk.ErrInternal("Address not authorized to join the contract.").Result()
	}

	if obj.State != StateCreated {
		return sdk.ErrInternal("Contract already joined.").Result()
	}

	if !msg.BobAmount.IsEqual(obj.BobAmount) {
		return sdk.ErrInternal("Invalid amount for Bob.").Result()
	}

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, obj.BobAddress, sdk.Coins{obj.BobAmount})
	if err != nil {
		return sdk.ErrInsufficientCoins("Not enough coins.").Result()
	}

	// Lock the contract, from now on both parties will need to sign any transaction.
	obj.State = StateLocked
	obj.Balance = obj.Balance.Plus(obj.BobAmount)
	keeper.UpsertContract(ctx, obj)

	return sdk.Result{}
}
