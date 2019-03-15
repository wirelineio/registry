//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/registry/x/multisig/msgs"
)

// Handle MsgInitMultiSig.
func handleMsgInitMultiSig(ctx sdk.Context, keeper Keeper, msg msgs.MsgInitMultiSig) sdk.Result {

	if keeper.HasContract(ctx, msg.ID) {
		return sdk.ErrInternal("Contract by that ID already exists.").Result()
	}

	if msg.AliceAddress.Equals(msg.BobAddress) {
		return sdk.ErrInternal("Alice and Bob addresses can't be identical.").Result()
	}

	if !msg.AliceAmount.IsPositive() {
		return sdk.ErrInternal("Invalid amount for Alice.").Result()
	}

	if !msg.BobAmount.IsPositive() {
		return sdk.ErrInternal("Invalid amount for Bob.").Result()
	}

	if !msg.AliceAmount.SameDenomAs(msg.BobAmount) {
		return sdk.ErrInternal("Amount denomination mismatch.").Result()
	}

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.AliceAddress, sdk.Coins{msg.AliceAmount})
	if err != nil {
		return sdk.ErrInsufficientCoins("Not enough coins.").Result()
	}

	obj := Contract{
		ID:           msg.ID,
		State:        StateCreated,
		AliceAmount:  msg.AliceAmount,
		AliceAddress: msg.AliceAddress,
		BobAmount:    msg.BobAmount,
		BobAddress:   msg.BobAddress,
		Balance:      msg.AliceAmount,
	}

	keeper.UpsertContract(ctx, obj)

	return sdk.Result{}
}
