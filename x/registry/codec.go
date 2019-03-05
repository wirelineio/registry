//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetResource{}, "registry/SetResource", nil)
	cdc.RegisterConcrete(MsgDeleteResource{}, "registry/DeleteResource", nil)
	cdc.RegisterConcrete(MsgClearResources{}, "registry/ClearResources", nil)
}
