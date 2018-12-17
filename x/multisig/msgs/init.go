//
// Copyright 2018 Wireline, Inc.
//

package msgs

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgInitMultiSig defines a InitMultiSig message.
type MsgInitMultiSig struct {
	ID           string
	Amount       sdk.Coin
	AliceAddress sdk.AccAddress
	BobAddress   sdk.AccAddress
}

// NewMsgInitMultiSig is the constructor function for MsgInitMultiSig.
func NewMsgInitMultiSig(id string, amount sdk.Coin, aliceAddress sdk.AccAddress, bobAddress sdk.AccAddress) MsgInitMultiSig {
	return MsgInitMultiSig{
		ID:           id,
		Amount:       amount,
		AliceAddress: aliceAddress,
		BobAddress:   bobAddress,
	}
}

// Route Implements Msg.
func (msg MsgInitMultiSig) Route() string { return "multisig" }

// Type Implements Msg.
func (msg MsgInitMultiSig) Type() string { return "init_multisig" }

// ValidateBasic Implements Msg.
func (msg MsgInitMultiSig) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return sdk.ErrUnknownRequest("ID cannot be empty.")
	}

	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive.")
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
func (msg MsgInitMultiSig) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgInitMultiSig) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AliceAddress}
}
