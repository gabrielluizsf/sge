// ecsign: the NIST curves and Ed25519, keys made from fixed bytes so the
// output is the same each run, and signatures checked both ways.
package main

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	for _, c := range []ecdh.Curve{ecdh.P256(), ecdh.P384(), ecdh.P521(), ecdh.X25519()} {
		a := fixedKey(c, 1)
		b := fixedKey(c, 2)
		s, err := a.ECDH(b.PublicKey())
		fmt.Println(c, hex.EncodeToString(a.PublicKey().Bytes())[:24], hex.EncodeToString(s)[:24], err)
	}
	h := sha256.Sum256([]byte("sge"))
	for _, c := range []elliptic.Curve{elliptic.P256(), elliptic.P384(), elliptic.P521()} {
		k, err := ecdsa.GenerateKey(c, rand.Reader)
		if err != nil {
			panic(err)
		}
		sig, err := ecdsa.SignASN1(rand.Reader, k, h[:])
		fmt.Println(c.Params().Name, err, ecdsa.VerifyASN1(&k.PublicKey, h[:], sig), ecdsa.VerifyASN1(&k.PublicKey, h[1:], sig))
	}
	seed := make([]byte, ed25519.SeedSize)
	seed[0] = 7
	k := ed25519.NewKeyFromSeed(seed)
	sig := ed25519.Sign(k, []byte("sge"))
	fmt.Println("ed25519", hex.EncodeToString(sig)[:24], ed25519.Verify(k.Public().(ed25519.PublicKey), []byte("sge"), sig))
}

// fixedKey is the private key whose scalar is n, padded to the curve's size.
func fixedKey(c ecdh.Curve, n byte) *ecdh.PrivateKey {
	size := map[ecdh.Curve]int{ecdh.P256(): 32, ecdh.P384(): 48, ecdh.P521(): 66, ecdh.X25519(): 32}[c]
	b := make([]byte, size)
	b[size-1] = n
	b[size-2] = 0x5a
	k, err := c.NewPrivateKey(b)
	if err != nil {
		panic(err)
	}
	return k
}
