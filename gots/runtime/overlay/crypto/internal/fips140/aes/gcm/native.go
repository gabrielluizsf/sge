//go:build gots_overlay

// AES-GCM by the runtime's AEAD, where it has one (rt/aead.ts).

package gcm

import "crypto/internal/fips140/aes"

// aeadJob is an operation for the runtime's AEAD; the runtime reads its
// fields by name.
type aeadJob struct {
	cipher                    string
	key, nonce, text, ad, out []byte
}

// aeadSeal and aeadOpen are the runtime's under gots: 0 when it has no
// AEAD, 1 when done, and 2 when what Open was given is not authentic.
func aeadSeal(j *aeadJob) int { return 0 }

func aeadOpen(j *aeadJob) int { return 0 }

// nativeJob is g's operation, if the runtime's AEAD can do it: a 12-byte
// nonce and a 16-byte tag.
func nativeJob(g *GCM, nonce []byte) (*aeadJob, bool) {
	if len(nonce) != gcmStandardNonceSize || g.tagSize != gcmTagSize {
		return nil, false
	}
	return &aeadJob{cipher: "aes-gcm", key: aes.GotsKey(&g.cipher), nonce: nonce}, true
}
