// curvefield: the NIST curves' and Curve25519's arithmetic, which runs on
// the overlay's own fields under gots, against Go's: points from fixed
// scalars, near the order too, compressed and uncompressed encodings, and
// encodings that must be refused.
package main

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	for _, c := range []elliptic.Curve{elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521()} {
		curve(c)
	}
	for _, c := range []ecdh.Curve{ecdh.P256(), ecdh.P384(), ecdh.P521(), ecdh.X25519()} {
		exchange(c)
	}
	signatures()
}

// curve multiplies fixed scalars, encodes, decodes and refuses points.
func curve(c elliptic.Curve) {
	p := c.Params()
	sum := sha256.New()
	nMinus1 := new(big.Int).Sub(p.N, big.NewInt(1))
	for _, k := range []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(0xdeadbeef), nMinus1, new(big.Int).Rsh(p.N, 1)} {
		x, y := c.ScalarBaseMult(k.Bytes())
		x2, y2 := c.ScalarMult(x, y, big.NewInt(12345).Bytes())
		x3, y3 := c.Add(x, y, x2, y2)
		x4, y4 := c.Double(x3, y3)
		for _, v := range []*big.Int{x, y, x2, y2, x3, y3, x4, y4} {
			sum.Write(v.Bytes())
		}
		comp := elliptic.MarshalCompressed(c, x4, y4)
		cx, cy := elliptic.UnmarshalCompressed(c, comp)
		ux, uy := elliptic.Unmarshal(c, elliptic.Marshal(c, x4, y4))
		sum.Write([]byte(fmt.Sprint(cx.Cmp(x4) == 0 && cy.Cmp(y4) == 0, ux.Cmp(x4) == 0 && uy.Cmp(y4) == 0)))
	}
	// An x of p and up, and a point off the curve, are refused.
	bad := elliptic.Marshal(c, p.Gx, p.Gy)
	size := (p.BitSize + 7) / 8
	pb := p.P.FillBytes(make([]byte, size))
	copy(bad[1:], pb)
	x, _ := elliptic.Unmarshal(c, bad)
	off := elliptic.Marshal(c, p.Gx, p.Gy)
	off[len(off)-1] ^= 1
	y, _ := elliptic.Unmarshal(c, off)
	fmt.Println(p.Name, hex.EncodeToString(sum.Sum(nil))[:32], x == nil, y == nil, c.IsOnCurve(p.Gx, p.Gy))
}

// exchange agrees on secrets from fixed keys, and refuses bad ones.
func exchange(c ecdh.Curve) {
	sum := sha256.New()
	size := map[ecdh.Curve]int{ecdh.P256(): 32, ecdh.P384(): 48, ecdh.P521(): 66, ecdh.X25519(): 32}[c]
	for i := 1; i < 6; i++ {
		a, b := key(c, size, byte(i)), key(c, size, byte(i*31))
		s1, err1 := a.ECDH(b.PublicKey())
		s2, err2 := b.ECDH(a.PublicKey())
		sum.Write(s1)
		sum.Write(a.PublicKey().Bytes())
		fmt.Sprint(err1, err2)
		if string(s1) != string(s2) {
			fmt.Println("disagree", c)
		}
	}
	_, err := c.NewPublicKey(make([]byte, size+1))
	fmt.Println(c, hex.EncodeToString(sum.Sum(nil))[:32], err != nil)
	if c == ecdh.X25519() {
		// A point of low order makes an all-zero secret, refused.
		low, _ := c.NewPublicKey(make([]byte, 32))
		_, err := key(c, size, 3).ECDH(low)
		fmt.Println("low order:", err)
	}
}

func key(c ecdh.Curve, size int, n byte) *ecdh.PrivateKey {
	b := make([]byte, size)
	for i := range b {
		b[i] = n ^ byte(i*13)
	}
	b[0] &= 0x7f
	if c == ecdh.P521() {
		b[0] = 1
	}
	k, err := c.NewPrivateKey(b)
	if err != nil {
		panic(err)
	}
	return k
}

// signatures checks ECDSA and Ed25519 signatures, good and tampered.
func signatures() {
	h := sha256.Sum256([]byte("curvefield"))
	for _, c := range []elliptic.Curve{elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521()} {
		d := new(big.Int).SetBytes(h[:])
		d.Mod(d, c.Params().N)
		k := &ecdsa.PrivateKey{D: d}
		k.Curve = c
		k.X, k.Y = c.ScalarBaseMult(d.Bytes())
		r, s, err := ecdsa.Sign(zeros{}, k, h[:])
		ok := ecdsa.Verify(&k.PublicKey, h[:], r, s)
		bad := ecdsa.Verify(&k.PublicKey, h[1:], r, s)
		fmt.Println(c.Params().Name, err, ok, bad)
	}
	for i := range 4 {
		seed := make([]byte, ed25519.SeedSize)
		seed[i] = byte(i + 1)
		k := ed25519.NewKeyFromSeed(seed)
		msg := []byte(fmt.Sprint("message", i))
		sig := ed25519.Sign(k, msg)
		pub := k.Public().(ed25519.PublicKey)
		tampered := append([]byte{}, sig...)
		tampered[i] ^= 4
		fmt.Println("ed25519", hex.EncodeToString(sig)[:24], hex.EncodeToString(pub)[:16],
			ed25519.Verify(pub, msg, sig), ed25519.Verify(pub, msg, tampered))
	}
}

// zeros is a reader of zero bytes, for deterministic signatures.
type zeros struct{}

func (zeros) Read(b []byte) (int, error) {
	clear(b)
	return len(b), nil
}
