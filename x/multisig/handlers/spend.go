//
// Copyright 2018 Wireline, Inc.
//

package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wirechain/x/multisig/msgs"
)

// Handle MsgSpendMultiSig.
func handleMsgSpendMultiSig(ctx sdk.Context, keeper Keeper, msg msgs.MsgSpendMultiSig) sdk.Result {
	if !keeper.HasContract(ctx, msg.ID) {
		return sdk.ErrInternal("Contract not found.").Result()
	}

	obj := keeper.GetContract(ctx, msg.ID)
	if !obj.AliceAddress.Equals(msg.AliceAddress) || !obj.BobAddress.Equals(msg.BobAddress) {
		return sdk.ErrInternal("Address not authorized to spend.").Result()
	}

	if obj.State != StateLocked {
		return sdk.ErrInternal("Contract counterparty hasn't joined, can't spend.").Result()
	}

	if obj.Balance.IsLT(msg.Amount) {
		return sdk.ErrInsufficientCoins("Not enough coins.").Result()
	}

	_, _, err := keeper.coinKeeper.AddCoins(ctx, msg.ToAddress, sdk.Coins{msg.Amount})
	if err != nil {
		return sdk.ErrInternal("Error transferring coins.").Result()
	}

	obj.Balance = obj.Balance.Minus(msg.Amount)
	keeper.UpsertContract(ctx, obj)

	return sdk.Result{}
}
