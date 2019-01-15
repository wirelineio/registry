//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccOutput represents an account based output birth record.
type AccOutput struct {
	ID      []byte
	Value   uint64
	Address sdk.AccAddress
	Block   int64
}

// OutPoint identifies an output from a previous transaction by index.
// Index >= 0 indicates Hash is a transaction ID.
// Index = -1 indicates Hash refers to an account based output birth record.
// Index = -2 indicates Hash refers to a voucher based output birth record.
type OutPoint struct {
	Hash  []byte
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
