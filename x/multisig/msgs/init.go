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
	AliceAmount  sdk.Coin
	AliceAddress sdk.AccAddress
	BobAmount    sdk.Coin
	BobAddress   sdk.AccAddress
}

// NewMsgInitMultiSig is the constructor function for MsgInitMultiSig.
func NewMsgInitMultiSig(id string, aliceAmount sdk.Coin, aliceAddress sdk.AccAddress, bobAmount sdk.Coin, bobAddress sdk.AccAddress) MsgInitMultiSig {
	return MsgInitMultiSig{
		ID:           id,
		AliceAmount:  aliceAmount,
		AliceAddress: aliceAddress,
		BobAmount:    bobAmount,
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

	if !msg.AliceAmount.IsPositive() {
		return sdk.ErrInsufficientCoins("Alice's amount must be positive.")
	}

	if msg.AliceAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.AliceAddress.String())
	}

	if !msg.BobAmount.IsPositive() {
		return sdk.ErrInsufficientCoins("Bob's amount must be positive.")
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
