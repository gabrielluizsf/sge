package main

import (
	"fmt"
	"math/big"
	"strings"
)

const limbBits = 22

// prime is what a template or the arithmetic needs of a spec's prime.
type prime struct {
	spec
	Type, Desc, Pre string
	Len, N          int // bytes, limbs
	P               *big.Int
	Limbs           []int64 // P's
	PInv            int64   // -p⁻¹ mod 2^limbBits
	PBytes          string  // P big-endian, as Go bytes
	Fields          string  // the element's
	Indices         []int   // 0 to N-1
}

func newPrime(s spec) *prime {
	n, _ := new(big.Int).SetString(s.hex, 16)
	p := &prime{spec: s, Type: s.typ, Desc: s.desc, Pre: s.pre, P: n, Len: (n.BitLen() + 7) / 8}
	p.N = (n.BitLen() + limbBits - 1) / limbBits
	p.Limbs = limbsOf(n, p.N)
	b := big.NewInt(1 << limbBits)
	inv := new(big.Int).ModInverse(new(big.Int).Mod(n, b), b)
	p.PInv = new(big.Int).Mod(inv.Neg(inv), b).Int64()
	var pb, fs []string
	for _, x := range n.FillBytes(make([]byte, p.Len)) {
		pb = append(pb, fmt.Sprintf("0x%02x", x))
	}
	for i := range p.N {
		fs = append(fs, fmt.Sprintf("l%d", i))
		p.Indices = append(p.Indices, i)
	}
	p.PBytes, p.Fields = strings.Join(pb, ", "), strings.Join(fs, ", ")
	return p
}

// limbsOf is x's n limbs, least significant first.
func limbsOf(x *big.Int, n int) []int64 {
	mask := big.NewInt(1<<limbBits - 1)
	var l []int64
	for i := range n {
		v := new(big.Int).Rsh(x, uint(i*limbBits))
		l = append(l, v.And(v, mask).Int64())
	}
	return l
}

// file is the path of the prime's file with suffix.
func (p *prime) file(suffix string) string { return p.dir + "/" + p.name + suffix }

// fn is the name of the prime's function op.
func (p *prime) fn(op string) string { return p.pre + op }
