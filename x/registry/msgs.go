//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSetRecord defines a SetResource message.
type MsgSetRecord struct {
	Payload PayloadObj
	Signer  sdk.AccAddress
}

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload PayloadObj, signer sdk.AccAddress) MsgSetRecord {
	return MsgSetRecord{
		Payload: payload,
		Signer:  signer,
	}
}

// Route Implements Msg.
func (msg MsgSetRecord) Route() string { return "registry" }

// Type Implements Msg.
func (msg MsgSetRecord) Type() string { return "set" }

// ValidateBasic Implements Msg.
func (msg MsgSetRecord) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	owner := msg.Payload.Record.Owner
	if owner == "" {
		return sdk.ErrInternal("Record owner not set.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetRecord) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSetRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgDeleteRecord defines a DeleteResource message.
type MsgDeleteRecord struct {
	Payload PayloadObj
	Signer  sdk.AccAddress
}

// NewMsgDeleteRecord is the constructor function for MsgDeleteRecord.
func NewMsgDeleteRecord(payload PayloadObj, signer sdk.AccAddress) MsgDeleteRecord {
	return MsgDeleteRecord{
		Payload: payload,
		Signer:  signer,
	}
}

// Route Implements Msg.
func (msg MsgDeleteRecord) Route() string { return "registry" }

// Type Implements Msg.
func (msg MsgDeleteRecord) Type() string { return "delete" }

// ValidateBasic Implements Msg.
func (msg MsgDeleteRecord) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	owner := msg.Payload.Record.Owner
	if owner == "" {
		return sdk.ErrInternal("Record owner not set.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDeleteRecord) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgDeleteRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgClearRecords defines a MsgClearRecords message.
type MsgClearRecords struct {
	Signer sdk.AccAddress
}

// NewMsgClearRecords is the constructor function for MsgClearRecords.
func NewMsgClearRecords(signer sdk.AccAddress) MsgClearRecords {
	return MsgClearRecords{
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgClearRecords) Route() string { return "registry" }

// Type Implements Msg.
func (msg MsgClearRecords) Type() string { return "clear" }

// ValidateBasic Implements Msg.
func (msg MsgClearRecords) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClearRecords) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgClearRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
