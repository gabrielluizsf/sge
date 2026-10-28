//go:build gots_overlay

// The key of a Block, for gcm's overlay to hand the runtime's AES-GCM:
// the expanded key starts with it.

package aes

import "crypto/internal/fips140deps/byteorder"

// GotsKey is the key c was made with.
func GotsKey(c *Block) []byte {
	key := make([]byte, 4*(c.rounds-6))
	for i := range len(key) / 4 {
		byteorder.BEPutUint32(key[4*i:], c.enc[i])
	}
	return key
}
