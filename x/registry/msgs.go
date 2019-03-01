//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSetResource defines a SetResource message.
type MsgSetResource struct {
	Payload PayloadObj
	Signer  sdk.AccAddress
}

// NewMsgSetResource is the constructor function for MsgSetResource.
func NewMsgSetResource(payload PayloadObj, signer sdk.AccAddress) MsgSetResource {
	return MsgSetResource{
		Payload: payload,
		Signer:  signer,
	}
}

// Route Implements Msg.
func (msg MsgSetResource) Route() string { return "registry" }

// Type Implements Msg.
func (msg MsgSetResource) Type() string { return "set" }

// ValidateBasic Implements Msg.
func (msg MsgSetResource) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	owner := msg.Payload.Resource.Owner
	if owner.Address == "" && owner.ID == "" {
		return sdk.ErrInternal("Resource owner not set.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetResource) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSetResource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
