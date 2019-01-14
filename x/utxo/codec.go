//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBirthAccOutput{}, "utxo/BirthAccOutput", nil)
}
