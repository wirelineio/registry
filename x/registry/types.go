//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// ID for resources.
type ID string

// Owner represents a resource owner.
type Owner struct {
	// If ID is populated, that will be used (ID of Owner resource record). Else, Address will be used.
	// One of the two MUST be populated.
	ID
	Address sdk.AccAddress
}

// Resource represents a registry record.
type Resource struct {
	ID
	Type             string
	Owner            Owner
	SystemAttributes []byte
	Attributes       []byte
}

// Signature represents a resource signature.
type Signature struct {
	crypto.PubKey
	Signature []byte
}

// Payload represents a signed resource payload.
type Payload struct {
	Resource
	Signatures []Signature
}
