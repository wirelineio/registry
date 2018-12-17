//
// Copyright 2018 Wireline, Inc.
//

package msgs

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSpendMultiSig defines a SpendMultiSig message.
type MsgSpendMultiSig struct {
	ID           string
	Amount       sdk.Coin
	ToAddress    sdk.AccAddress
	AliceAddress sdk.AccAddress
	BobAddress   sdk.AccAddress
}

// NewMsgSpendMultiSig is the constructor function for MsgSpendMultiSig.
func NewMsgSpendMultiSig(id string, amount sdk.Coin, toAddress sdk.AccAddress, aliceAddress sdk.AccAddress, bobAddress sdk.AccAddress) MsgSpendMultiSig {
	return MsgSpendMultiSig{
		ID:           id,
		Amount:       amount,
		ToAddress:    toAddress,
		AliceAddress: aliceAddress,
		BobAddress:   bobAddress,
	}
}

// Route Implements Msg.
func (msg MsgSpendMultiSig) Route() string { return "multisig" }

// Type Implements Msg.
func (msg MsgSpendMultiSig) Type() string { return "spend_multisig" }

// ValidateBasic Implements Msg.
func (msg MsgSpendMultiSig) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return sdk.ErrUnknownRequest("ID cannot be empty.")
	}

	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive.")
	}

	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.ToAddress.String())
	}

	if msg.AliceAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.AliceAddress.String())
	}

	if msg.BobAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.BobAddress.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSpendMultiSig) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSpendMultiSig) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AliceAddress, msg.BobAddress}
}
