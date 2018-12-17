//
// Copyright 2018 Wireline, Inc.
//

package multisig

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/wirelineio/wirechain/x/multisig/msgs"
)

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(msgs.MsgInitMultiSig{}, "multisig/InitMultiSig", nil)
	cdc.RegisterConcrete(msgs.MsgAbortMultiSig{}, "multisig/AbortMultiSig", nil)
	cdc.RegisterConcrete(msgs.MsgJoinMultiSig{}, "multisig/JoinMultiSig", nil)
	cdc.RegisterConcrete(msgs.MsgSpendMultiSig{}, "multisig/SpendMultiSig", nil)
}
