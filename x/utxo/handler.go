//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "utxo" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgBirthAccOutput:
			return handleMsgBirthAccOutput(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized utxo Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgBirthAccOutput.
func handleMsgBirthAccOutput(ctx sdk.Context, keeper Keeper, msg MsgBirthAccOutput) sdk.Result {

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Address, sdk.Coins{msg.Amount})
	if err != nil {
		return sdk.ErrInsufficientCoins("Not enough coins to create UTXO.").Result()
	}

	// Create AccOutput record.
	accUtxo, err := GenAccOutput(ctx, keeper, msg)
	if err != nil {
		return sdk.ErrInternal("Error generating account UTXO.").Result()
	}

	keeper.PutAccOutput(ctx, accUtxo)
	keeper.PutOutPoint(ctx, OutPoint{
		Hash:  accUtxo.ID,
		Index: OutPointAccountBirth,
	})

	return sdk.Result{}
}
