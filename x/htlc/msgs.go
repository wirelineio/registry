package htlc

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Status represents the status of an HTLC.
type Status uint8

// HTLC status enum.
const (
	HtlcCreated  Status = 1
	HtlcRedeemed Status = 2
	HtlcFailed   Status = 3
)

// MsgAddHtlc defines a AddHtlc message.
type MsgAddHtlc struct {
	Status         Status
	Amount         sdk.Coin
	Hash           string
	Timelock       uint
	RedeemAddress  sdk.AccAddress
	TimeoutAddress sdk.AccAddress
}

// NewMsgAddHtlc is the constructor function for MsgAddHtlc.
func NewMsgAddHtlc(amount sdk.Coin, hash string, timelock uint, redeemAddress sdk.AccAddress, timeoutAddress sdk.AccAddress) MsgAddHtlc {
	return MsgAddHtlc{
		Status:         HtlcCreated,
		Amount:         amount,
		Hash:           hash,
		Timelock:       timelock,
		RedeemAddress:  redeemAddress,
		TimeoutAddress: timeoutAddress,
	}
}

// Route Implements Msg.
func (msg MsgAddHtlc) Route() string { return "htlc" }

// Type Implements Msg.
func (msg MsgAddHtlc) Type() string { return "add_htlc" }

// ValidateBasic Implements Msg.
func (msg MsgAddHtlc) ValidateBasic() sdk.Error {
	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive.")
	}

	if len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("Hash cannot be empty.")
	}

	if msg.Timelock == 0 {
		return sdk.ErrUnknownRequest("Timelock should be greater than zero.")
	}

	if msg.RedeemAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.RedeemAddress.String())
	}

	if msg.TimeoutAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.RedeemAddress.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddHtlc) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgAddHtlc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.TimeoutAddress}
}

// MsgRedeemHtlc defines the RedeemHtlc message.
type MsgRedeemHtlc struct {
	Preimage string
	Redeemer sdk.AccAddress
}

// NewMsgRedeemHtlc is the constructor function for MsgRedeemHtlc.
func NewMsgRedeemHtlc(preimage string, redeemer sdk.AccAddress) MsgRedeemHtlc {
	return MsgRedeemHtlc{
		Preimage: preimage,
		Redeemer: redeemer,
	}
}

// Route Implements Msg.
func (msg MsgRedeemHtlc) Route() string { return "htlc" }

// Type Implements Msg.
func (msg MsgRedeemHtlc) Type() string { return "redeem_htlc" }

// ValidateBasic Implements Msg.
func (msg MsgRedeemHtlc) ValidateBasic() sdk.Error {
	if len(msg.Preimage) == 0 {
		return sdk.ErrUnknownRequest("Preimage cannot be empty.")
	}

	if msg.Redeemer.Empty() {
		return sdk.ErrInvalidAddress(msg.Redeemer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRedeemHtlc) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgRedeemHtlc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Redeemer}
}

// MsgFailHtlc defines the FailHtlc message.
type MsgFailHtlc struct {
	Hash   string
	Sender sdk.AccAddress
}

// NewMsgFailHtlc is the constructor function for MsgFailHtlc.
func NewMsgFailHtlc(hash string, sender sdk.AccAddress) MsgFailHtlc {
	return MsgFailHtlc{
		Hash:   hash,
		Sender: sender,
	}
}

// Route Implements Msg.
func (msg MsgFailHtlc) Route() string { return "htlc" }

// Type Implements Msg.
func (msg MsgFailHtlc) Type() string { return "fail_htlc" }

// ValidateBasic Implements Msg.
func (msg MsgFailHtlc) ValidateBasic() sdk.Error {
	if len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("Hash cannot be empty.")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgFailHtlc) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgFailHtlc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
