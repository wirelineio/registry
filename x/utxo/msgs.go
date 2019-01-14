//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"encoding/json"

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
