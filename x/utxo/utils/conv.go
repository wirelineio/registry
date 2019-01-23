//
// Copyright 2019 Wireline, Inc.
//

package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// UInt64ToBytes converts a unint64 into a byte array.
func UInt64ToBytes(n uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, n)

	return buf
}

// Int64ToBytes converts a unint64 into a byte array.
func Int64ToBytes(n int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, n)

	return buf
}

// BytesToBase64 encodes a byte array as a base64 string.
func BytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// BytesToHex encodes a byte array as an upper case hex string.
func BytesToHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}
