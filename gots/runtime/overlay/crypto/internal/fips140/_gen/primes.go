package main

// spec is what tells a prime's generated code apart.
type spec struct {
	dir, pkg  string // where it goes
	typ, pre  string // the element's type, and its functions' prefix
	name      string // the file's name, without .go
	desc, hex string // the prime, as Go's docs write it, and in hex
}

var nistPrimes = []spec{
	nist("224", "2^224 - 2^96 + 1", "ffffffffffffffffffffffffffffffff000000000000000000000001"),
	nist("256", "2^256 - 2^224 + 2^192 + 2^96 - 1", "ffffffff00000001000000000000000000000000ffffffffffffffffffffffff"),
	nist("384", "2^384 - 2^128 - 2^96 + 2^32 - 1", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffff0000000000000000ffffffff"),
	nist("521", "2^521 - 1", "01ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
}

func nist(bits, desc, hex string) spec {
	return spec{dir: "nistec/fiat", pkg: "fiat", typ: "P" + bits + "Element", pre: "p" + bits,
		name: "p" + bits, desc: desc, hex: hex}
}

var field25519 = spec{dir: "edwards25519/field", pkg: "field", typ: "Element", pre: "fe",
	name: "fe", desc: "2^255 - 19", hex: "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed"}
