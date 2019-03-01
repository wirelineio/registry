//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/tendermint/tendermint/crypto"
)

// GetResourceSignature returns a cryptographic signature for a transaction.
func GetResourceSignature(resource Resource, name string) ([]byte, crypto.PubKey, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, nil, err
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return nil, nil, err
	}

	signBytes := GenResourceHash(resource)

	sigBytes, pubKey, err := keybase.Sign(name, passphrase, signBytes)
	if err != nil {
		return nil, nil, err
	}

	return sigBytes, pubKey, nil
}
