//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"bytes"
	"encoding/json"
	"sort"

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

// NewMsgTx creates a new transaction message.
func NewMsgTx(tx Tx, signer sdk.AccAddress) MsgTx {
	return MsgTx{
		Signer: signer,
		Tx:     tx,
	}
}

// NewTxPayToAddress creates a transaction payload to pay to an address.
func NewTxPayToAddress(cdc *codec.Codec, sig []byte, hash []byte, index int32, amount uint64, change uint64, from sdk.AccAddress, to sdk.AccAddress) Tx {
	tx := Tx{
		TxIn: []TxIn{
			TxIn{
				Input: OutPoint{
					Hash:  hash,
					Index: index,
				},
				Witness: sig,
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
	}

	SortTxInputs(&tx)
	SortTxOutputs(&tx)

	return tx
}

// SortTxInputs sorts transaction inputs (canonical ordering).
func SortTxInputs(tx *Tx) {
	sort.SliceStable(tx.TxIn, func(i, j int) bool {
		a := tx.TxIn[i]
		b := tx.TxIn[i]

		bytesCompare := bytes.Compare([]byte(a.Input.Hash), []byte(b.Input.Hash))

		if (bytesCompare < 0) || (bytesCompare == 0 && (a.Input.Index < b.Input.Index)) {
			return true
		}

		return false
	})
}

// SortTxOutputs sorts transaction outputs (canonical ordering).
func SortTxOutputs(tx *Tx) {
	sort.SliceStable(tx.TxOut, func(i, j int) bool {
		a := tx.TxOut[i]
		b := tx.TxOut[j]

		if a.Value < b.Value {
			return true
		}

		if a.Value == b.Value {
			bytesCompare := bytes.Compare([]byte(a.PkScript), []byte(b.PkScript))
			if bytesCompare < 0 {
				return true
			}
		}

		return false
	})
}
