//
// Copyright 2018 Wireline, Inc.
//

package msgs

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgJoinMultiSig defines a JoinMultiSig message.
type MsgJoinMultiSig struct {
	ID         string
	BobAmount  sdk.Coin
	BobAddress sdk.AccAddress
}

// NewMsgJoinMultiSig is the constructor function for MsgJoinMultiSig.
func NewMsgJoinMultiSig(id string, bobAmount sdk.Coin, bobAddress sdk.AccAddress) MsgJoinMultiSig {
	return MsgJoinMultiSig{
		ID:         id,
		BobAmount:  bobAmount,
		BobAddress: bobAddress,
	}
}

// Route Implements Msg.
func (msg MsgJoinMultiSig) Route() string { return "multisig" }

// Type Implements Msg.
func (msg MsgJoinMultiSig) Type() string { return "join_multisig" }

// ValidateBasic Implements Msg.
func (msg MsgJoinMultiSig) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return sdk.ErrUnknownRequest("ID cannot be empty.")
	}

	if !msg.BobAmount.IsPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive.")
	}

	if msg.BobAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.BobAddress.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgJoinMultiSig) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgJoinMultiSig) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.BobAddress}
}
