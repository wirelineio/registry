//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetRecord{}, "registry/SetResource", nil)
	cdc.RegisterConcrete(MsgDeleteRecord{}, "registry/DeleteResource", nil)
	cdc.RegisterConcrete(MsgClearRecords{}, "registry/ClearResources", nil)
}
