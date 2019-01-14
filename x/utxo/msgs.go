//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBirthAccUtxo defines a BirthAccUtxo message.
type MsgBirthAccUtxo struct {
	Amount  sdk.Coin
	Address sdk.AccAddress
}

// NewMsgBirthAccUtxo is the constructor function for MsgBirthAccUtxo.
func NewMsgBirthAccUtxo(amount sdk.Coin, address sdk.AccAddress) MsgBirthAccUtxo {
	return MsgBirthAccUtxo{
		Amount:  amount,
		Address: address,
	}
}

// Route Implements Msg.
func (msg MsgBirthAccUtxo) Route() string { return "utxo" }

// Type Implements Msg.
func (msg MsgBirthAccUtxo) Type() string { return "new_utxo" }

// ValidateBasic Implements Msg.
func (msg MsgBirthAccUtxo) ValidateBasic() sdk.Error {
	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive.")
	}

	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBirthAccUtxo) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgBirthAccUtxo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
