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

// MsgDeleteResource defines a DeleteResource message.
type MsgDeleteResource struct {
	Payload PayloadObj
	Signer  sdk.AccAddress
}

// NewMsgDeleteResource is the constructor function for MsgDeleteResource.
func NewMsgDeleteResource(payload PayloadObj, signer sdk.AccAddress) MsgDeleteResource {
	return MsgDeleteResource{
		Payload: payload,
		Signer:  signer,
	}
}

// Route Implements Msg.
func (msg MsgDeleteResource) Route() string { return "registry" }

// Type Implements Msg.
func (msg MsgDeleteResource) Type() string { return "delete" }

// ValidateBasic Implements Msg.
func (msg MsgDeleteResource) ValidateBasic() sdk.Error {

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
func (msg MsgDeleteResource) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgDeleteResource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgClearResources defines a MsgClearResources message.
type MsgClearResources struct {
	Signer sdk.AccAddress
}

// NewMsgClearResources is the constructor function for MsgClearResources.
func NewMsgClearResources(signer sdk.AccAddress) MsgClearResources {
	return MsgClearResources{
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgClearResources) Route() string { return "registry" }

// Type Implements Msg.
func (msg MsgClearResources) Type() string { return "clear" }

// ValidateBasic Implements Msg.
func (msg MsgClearResources) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClearResources) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgClearResources) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
