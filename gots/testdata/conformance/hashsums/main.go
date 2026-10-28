// hashsums: the SHA-2 and SHA-3 sums, SHAKE, ML-KEM and HMAC of inputs of every length across a
// few blocks, and AES-GCM, whose block functions and GHASH run as the
// runtime's under gots.
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/mlkem"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

func main() {
	data := make([]byte, 300)
	for i := range data {
		data[i] = byte(i*7 + 3)
	}
	sha3New := func() hash.Hash { return sha3.New256() }
	for _, h := range []func() hash.Hash{sha256.New, sha512.New, sha512.New384, sha512.New512_224, sha3New} {
		all := h()
		for n := range data {
			d := h()
			d.Write(data[:n])
			all.Write(d.Sum(nil))
		}
		fmt.Println(hex.EncodeToString(all.Sum(nil)))
	}
	shake := sha3.NewSHAKE128()
	shake.Write(data)
	out := make([]byte, 600)
	shake.Read(out)
	fmt.Println(hex.EncodeToString(out[560:]))
	seed := make([]byte, mlkem.SeedSize)
	copy(seed, data)
	dk, _ := mlkem.NewDecapsulationKey768(seed)
	ek := dk.EncapsulationKey().Bytes()
	fmt.Println(hex.EncodeToString(ek[:16]), hex.EncodeToString(ek[len(ek)-16:]))
	shared, ct := dk.EncapsulationKey().Encapsulate()
	got, err := dk.Decapsulate(ct)
	fmt.Println(string(got) == string(shared), err)
	m := hmac.New(sha512.New, data[:40])
	m.Write(data)
	fmt.Println(hex.EncodeToString(m.Sum(nil)))
	block, _ := aes.NewCipher(data[:32])
	gcm, _ := cipher.NewGCM(block)
	for _, n := range []int{0, 1, 15, 16, 17, 100, 300} {
		sealed := gcm.Seal(nil, data[:12], data[:n], data[n%7:n%7+n%20])
		opened, err := gcm.Open(nil, data[:12], sealed, data[n%7:n%7+n%20])
		ad := data[n%7 : n%7+n%20]
		// Tampering with the text, the tag, the data or the nonce is
		// caught (CVE-2025-24015: some Deno did not check the tag).
		badText := tamper(sealed, 0)
		badTag := tamper(sealed, len(sealed)-1)
		_, e1 := gcm.Open(nil, data[:12], badText, ad)
		_, e2 := gcm.Open(nil, data[:12], badTag, ad)
		_, e3 := gcm.Open(nil, data[:12], sealed, append([]byte{1}, ad...))
		_, e4 := gcm.Open(nil, data[1:13], sealed, ad)
		fmt.Println(n, hex.EncodeToString(sealed[max(0, len(sealed)-20):]), len(opened) == n, err, e1 != nil, e2 != nil, e3 != nil, e4 != nil)
	}
}

// tamper is b with the byte at i flipped.
func tamper(b []byte, i int) []byte {
	t := append([]byte{}, b...)
	t[i] ^= 1
	return t
}
