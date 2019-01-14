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
		case MsgBirthAccUtxo:
			return handleMsgBirthAccUtxo(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized utxo Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgBirthAccUtxo.
func handleMsgBirthAccUtxo(ctx sdk.Context, keeper Keeper, msg MsgBirthAccUtxo) sdk.Result {

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Address, sdk.Coins{msg.Amount})
	if err != nil {
		return sdk.ErrInsufficientCoins("Not enough coins to create UTXO.").Result()
	}

	// Create AccountUtxo record.
	accUtxo, err := GenAccountUtxo(ctx, keeper, msg)
	if err != nil {
		return sdk.ErrInternal("Error generating account UTXO.").Result()
	}

	keeper.PutAccountUtxo(ctx, accUtxo)

	return sdk.Result{}
}
