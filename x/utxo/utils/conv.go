//
// Copyright 2019 Wireline, Inc.
//

package utils

import "encoding/binary"

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
