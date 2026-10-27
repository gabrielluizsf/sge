//go:build gots_overlay

// Overlay for sealGeneric and openGeneric, which hand the runtime's AEAD
// what it can do (see native.go), and do the rest as Go does.

package gcm

import (
	"crypto/internal/fips140/aes"
	"crypto/internal/fips140/subtle"
)

func sealGeneric(out []byte, g *GCM, nonce, plaintext, additionalData []byte) {
	if j, ok := nativeJob(g, nonce); ok {
		j.text, j.ad, j.out = plaintext, additionalData, out
		if aeadSeal(j) == 1 {
			return
		}
	}
	var H, counter, tagMask [gcmBlockSize]byte
	aes.EncryptBlockInternal(&g.cipher, H[:], H[:])
	deriveCounterGeneric(&H, &counter, nonce)
	gcmCounterCryptGeneric(&g.cipher, tagMask[:], tagMask[:], &counter)

	gcmCounterCryptGeneric(&g.cipher, out, plaintext, &counter)

	var tag [gcmTagSize]byte
	gcmAuthGeneric(tag[:], &H, &tagMask, out[:len(plaintext)], additionalData)
	copy(out[len(plaintext):], tag[:])
}

func openGeneric(out []byte, g *GCM, nonce, ciphertext, additionalData []byte) error {
	if j, ok := nativeJob(g, nonce); ok {
		j.text, j.ad, j.out = ciphertext, additionalData, out
		switch aeadOpen(j) {
		case 1:
			return nil
		case 2:
			return errOpen
		}
	}
	var H, counter, tagMask [gcmBlockSize]byte
	aes.EncryptBlockInternal(&g.cipher, H[:], H[:])
	deriveCounterGeneric(&H, &counter, nonce)
	gcmCounterCryptGeneric(&g.cipher, tagMask[:], tagMask[:], &counter)

	tag := ciphertext[len(ciphertext)-g.tagSize:]
	ciphertext = ciphertext[:len(ciphertext)-g.tagSize]

	var expectedTag [gcmTagSize]byte
	gcmAuthGeneric(expectedTag[:], &H, &tagMask, ciphertext, additionalData)
	if subtle.ConstantTimeCompare(expectedTag[:g.tagSize], tag) != 1 {
		return errOpen
	}

	gcmCounterCryptGeneric(&g.cipher, out, ciphertext, &counter)

	return nil
}
