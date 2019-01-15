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
		case MsgTx:
			return handleMsgTx(ctx, keeper, msg)
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

// Handle MsgTx.
func handleMsgTx(ctx sdk.Context, keeper Keeper, msg MsgTx) sdk.Result {

	// Supports only 1 input, for now.
	if len(msg.Tx.TxIn) != 1 {
		return sdk.ErrInternal("Multiple inputs not yet supported.").Result()
	}

	input := msg.Tx.TxIn[0].Input

	// Check that the input outpoint is in the UTXO list.
	if !keeper.HasOutPoint(ctx, input) {
		return sdk.ErrUnauthorized("OutPoint not found or already spent.").Result()
	}

	var inputValue uint64
	outputValue := GetTxOutValue(msg.Tx.TxOut)

	var redeemAddress sdk.AccAddress

	if input.Index >= 0 {
		tx := keeper.GetTx(ctx, input.Hash)
		txOut := tx.TxOut[input.Index]
		var obj PayToAddress
		keeper.cdc.MustUnmarshalBinaryBare(txOut.PkScript, &obj)

		redeemAddress = obj.Address
		inputValue = txOut.Value
	} else if input.Index == OutPointAccountBirth {
		accOutput := keeper.GetAccOutput(ctx, input.Hash)

		redeemAddress = accOutput.Address
		inputValue = accOutput.Value
	}

	if inputValue != outputValue {
		fmt.Println("input output", inputValue, outputValue)
		return sdk.ErrUnauthorized("Mismatch between input and output values.").Result()
	}

	if !redeemAddress.Equals(msg.Signer) {
		return sdk.ErrUnauthorized("OutPoint not spendable by message signer.").Result()
	}

	// Save Tx.
	txHash := GenTxHash(keeper, msg.Tx)
	keeper.PutTx(ctx, txHash, msg.Tx)

	// Delete old UTXO.
	keeper.DeleteOutPoint(ctx, input)

	// Create new UTXOs.
	for index := range msg.Tx.TxOut {
		keeper.PutOutPoint(ctx, OutPoint{
			Hash:  txHash,
			Index: int32(index),
		})
	}

	return sdk.Result{}
}
