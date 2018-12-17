//
// Copyright 2018 Wireline, Inc.
//

package msgs

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgAbortMultiSig defines a AbortMultiSig message.
type MsgAbortMultiSig struct {
	ID           string
	AliceAddress sdk.AccAddress
}

// NewMsgAbortMultiSig is the constructor function for MsgAbortMultiSig.
func NewMsgAbortMultiSig(id string, aliceAddress sdk.AccAddress) MsgAbortMultiSig {
	return MsgAbortMultiSig{
		ID:           id,
		AliceAddress: aliceAddress,
	}
}

// Route Implements Msg.
func (msg MsgAbortMultiSig) Route() string { return "multisig" }

// Type Implements Msg.
func (msg MsgAbortMultiSig) Type() string { return "abort_multisig" }

// ValidateBasic Implements Msg.
func (msg MsgAbortMultiSig) ValidateBasic() sdk.Error {
	if len(msg.ID) == 0 {
		return sdk.ErrUnknownRequest("ID cannot be empty.")
	}

	if msg.AliceAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.AliceAddress.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAbortMultiSig) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgAbortMultiSig) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AliceAddress}
}
