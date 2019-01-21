//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Hash represents a transaction or account output ID.
type Hash []byte

// MarshalJSON marshals to JSON using Bech32.
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// String implements the Stringer interface.
func (h Hash) String() string {
	return strings.ToUpper(hex.EncodeToString(h))
}

// AccOutput represents an account based output birth record.
type AccOutput struct {
	ID      Hash
	Value   uint64
	Address sdk.AccAddress
	Block   int64
}

// OutPoint identifies an output from a previous transaction by index.
// Index >= 0 indicates Hash is a transaction ID.
// Index = -1 indicates Hash refers to an account based output birth record.
// Index = -2 indicates Hash refers to a voucher based output birth record.
type OutPoint struct {
	Hash  Hash
	Index int32
}

// OutPointAccountBirth indicates Hash refers to an account based output birth record.
const OutPointAccountBirth = -1

// PayToAddress indicates the UTXO is payable to an address.
type PayToAddress struct {
	Address sdk.AccAddress
}

// PayToScript indicates the UTXO is payable to a script.
// type PayToScript struct {
// 	Script []byte
// }

// TxOut represents a transaction output.
// PkScript is the go-amino binary marshalled PayTo* struct.
type TxOut struct {
	Value    uint64
	PkScript []byte
}

// TxIn represents a transaction input.
type TxIn struct {
	Input    OutPoint
	Witness  []byte
	Sequence uint32
}

// Tx represents a transaction.
type Tx struct {
	TxIn     []TxIn
	TxOut    []TxOut
	LockTime uint32
}

// Wallet data structures.

// OutPointVal is an outpoint and it's value.
type OutPointVal struct {
	Hash  Hash
	Index int32
	Value uint64
}

// Wallet represents a balance and UTXOs for an address.
type Wallet struct {
	Balance uint64
	Entries []OutPointVal
}
