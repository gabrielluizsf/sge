//go:build gots_overlay

// Overlay for seal and open, which hand the runtime's ChaCha20-Poly1305
// (rt/aead.ts) the work, with Go's checks of overlaps first, and do it as
// Go does where the runtime has none.

package chacha20poly1305

import "golang.org/x/crypto/internal/alias"

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

func (c *chacha20poly1305) seal(dst, nonce, plaintext, additionalData []byte) []byte {
	ret, out := sliceForAppend(dst, len(plaintext)+16)
	if alias.InexactOverlap(out, plaintext) {
		panic("chacha20poly1305: invalid buffer overlap of output and input")
	}
	if alias.AnyOverlap(out, additionalData) {
		panic("chacha20poly1305: invalid buffer overlap of output and additional data")
	}
	j := &aeadJob{cipher: "chacha20-poly1305", key: c.key[:], nonce: nonce, text: plaintext, ad: additionalData, out: out}
	if aeadSeal(j) == 1 {
		return ret
	}
	return c.sealGeneric(dst, nonce, plaintext, additionalData)
}

func (c *chacha20poly1305) open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	ret, out := sliceForAppend(dst, len(ciphertext)-16)
	if alias.InexactOverlap(out, ciphertext) {
		panic("chacha20poly1305: invalid buffer overlap of output and input")
	}
	if alias.AnyOverlap(out, additionalData) {
		panic("chacha20poly1305: invalid buffer overlap of output and additional data")
	}
	j := &aeadJob{cipher: "chacha20-poly1305", key: c.key[:], nonce: nonce, text: ciphertext, ad: additionalData, out: out}
	switch aeadOpen(j) {
	case 1:
		return ret, nil
	case 2:
		return nil, errOpen
	}
	return c.openGeneric(dst, nonce, ciphertext, additionalData)
}
