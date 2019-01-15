//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"crypto/sha256"

	"github.com/wirelineio/wirechain/x/utxo/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenAccOutput creates an AccOutput from MsgBirthAccOutput info.
func GenAccOutput(ctx sdk.Context, keeper Keeper, msg MsgBirthAccOutput) (AccOutput, sdk.Error) {

	sequence, err := keeper.accountKeeper.GetSequence(ctx, msg.Address)
	if err != nil {
		return AccOutput{}, err
	}

	hash := sha256.New()
	hash.Write(msg.Address)
	hash.Write([]byte(msg.Amount.String()))
	hash.Write(utils.UInt64ToBytes(sequence))
	hash.Write(utils.Int64ToBytes(ctx.BlockHeight()))
	id := hash.Sum(nil)

	return AccOutput{
		ID:      id,
		Value:   uint64(msg.Amount.Amount.Int64()),
		Address: msg.Address,
		Block:   ctx.BlockHeight(),
	}, nil
}

// GetTxOutValue returns the sum of the output values.
func GetTxOutValue(outputs []TxOut) uint64 {
	var value uint64

	for _, output := range outputs {
		value += output.Value
	}

	return value
}

// GenTxHash generates a transaction hash.
func GenTxHash(keeper Keeper, tx Tx) []byte {
	// TODO(ashwin): Sort inputs/outputs in canonical order.

	first := sha256.New()
	first.Write(keeper.cdc.MustMarshalBinaryBare(tx))
	firstHash := first.Sum(nil)

	second := sha256.New()
	second.Write(firstHash)
	secondHash := second.Sum(nil)

	return secondHash
}