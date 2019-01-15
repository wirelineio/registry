//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBirthAccOutput defines a BirthAccOutput message.
type MsgBirthAccOutput struct {
	Amount  sdk.Coin
	Address sdk.AccAddress
}

// NewMsgBirthAccOutput is the constructor function for MsgBirthAccOutput.
func NewMsgBirthAccOutput(amount sdk.Coin, address sdk.AccAddress) MsgBirthAccOutput {
	return MsgBirthAccOutput{
		Amount:  amount,
		Address: address,
	}
}

// Route Implements Msg.
func (msg MsgBirthAccOutput) Route() string { return "utxo" }

// Type Implements Msg.
func (msg MsgBirthAccOutput) Type() string { return "birth_acc_output" }

// ValidateBasic Implements Msg.
func (msg MsgBirthAccOutput) ValidateBasic() sdk.Error {
	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive.")
	}

	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBirthAccOutput) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgBirthAccOutput) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgTx represents a UTXO based transaction.
type MsgTx struct {
	Tx     Tx
	Signer sdk.AccAddress
}

// Route Implements Msg.
func (msg MsgTx) Route() string { return "utxo" }

// Type Implements Msg.
func (msg MsgTx) Type() string { return "tx" }

// ValidateBasic Implements Msg.
func (msg MsgTx) ValidateBasic() sdk.Error {
	if len(msg.Tx.TxIn) == 0 {
		return sdk.ErrInternal("Must have at least one input.")
	}

	if len(msg.Tx.TxOut) == 0 {
		return sdk.ErrInternal("Must have at least one output.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTx) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// NewMsgPayToAddress creates a transaction to pay to an address.
func NewMsgPayToAddress(cdc *codec.Codec, hash []byte, index int32, amount uint64, change uint64, from sdk.AccAddress, to sdk.AccAddress) MsgTx {
	// TODO(ashwin): Canonical sorting of inputs/outputs.
	return MsgTx{
		Signer: from,
		Tx: Tx{
			TxIn: []TxIn{
				TxIn{
					Input: OutPoint{
						Hash:  hash,
						Index: index,
					},
					Witness: from.Bytes(),
				},
			},
			TxOut: []TxOut{
				TxOut{
					Value: amount,
					PkScript: cdc.MustMarshalBinaryBare(PayToAddress{
						Address: to,
					}),
				},
				TxOut{
					Value: change,
					PkScript: cdc.MustMarshalBinaryBare(PayToAddress{
						Address: from,
					}),
				},
			},
		},
	}
}
