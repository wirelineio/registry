//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/tendermint/tendermint/crypto"
	"golang.org/x/crypto/ripemd160"
)

// GenResourceHash generates a transaction hash.
func GenResourceHash(r ResourceYaml) []byte {
	first := sha256.New()

	bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		panic("Resource marshal error.")
	}

	first.Write(bytes)
	firstHash := first.Sum(nil)

	second := sha256.New()
	second.Write(firstHash)
	secondHash := second.Sum(nil)

	return secondHash
}

// GetAddressFromPubKey gets an address from the public key.
func GetAddressFromPubKey(pubKey crypto.PubKey) string {
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(pubKey.Bytes())
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha)
	ripemd := hasherRIPEMD160.Sum(nil)

	return BytesToHex(ripemd)
}

// BytesToBase64 encodes a byte array as a base64 string.
func BytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// BytesFromBase64 decodes a byte array from a base64 string.
func BytesFromBase64(str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic("Error decoding string to bytes.")
	}

	return bytes
}

// BytesToHex encodes a byte array as a hex string.
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// BytesFromHex decodes a byte array from a hex string.
func BytesFromHex(str string) []byte {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic("Error decoding hex to bytes.")
	}

	return bytes
}
