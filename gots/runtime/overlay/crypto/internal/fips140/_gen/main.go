// Command _gen writes the overlay's field arithmetic under gots: for each
// prime, Montgomery multiplication, sums and differences unrolled on
// 22-bit limbs (see nistec/fiat/mont.go), with the prime's limbs as
// constants, so a product by a zero limb is not there at all. It writes
// fiat's elements whole, and edwards25519/field's arithmetic, whose API
// is written by hand. Run it from the fips140 overlay's directory:
//
//	go run ./_gen
package main

func main() {
	for _, s := range nistPrimes {
		p := newPrime(s)
		write(p.file(".go"), element(p))
		write(p.file("_arith.go"), arith(p))
	}
	p := newPrime(field25519)
	write(p.file("_arith.go"), arith(p))
	write(p.file("_consts.go"), consts(p))
}
