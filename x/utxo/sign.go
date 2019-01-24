//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetTxSignature returns a cryptographic signature for a transaction.
func GetTxSignature(cdc *codec.Codec, tx Tx, name string) ([]byte, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return nil, err
	}

	txHash := GenTxHash(cdc, tx)

	sigBytes, _, err := keybase.Sign(name, passphrase, txHash)
	if err != nil {
		return nil, err
	}

	return sigBytes, nil
}
