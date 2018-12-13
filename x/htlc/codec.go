//
// Copyright 2018 Wireline, Inc.
//

package htlc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddHtlc{}, "htlc/AddHtlc", nil)
	cdc.RegisterConcrete(MsgRedeemHtlc{}, "htlc/RedeemHtlc", nil)
	cdc.RegisterConcrete(MsgFailHtlc{}, "htlc/FailHtlc", nil)
	cdc.RegisterConcrete(MsgClearHtlc{}, "htlc/ClearHtlc", nil)
}
