package main

import (
	"fmt"
	"math/big"
	"strings"
)

// consts writes the Montgomery constants of p's elements, R and R² mod p,
// as element values, for a package with no montField of its own.
func consts(p *prime) []byte {
	var o out
	o.header(p)
	r := new(big.Int).Lsh(big.NewInt(1), uint(limbBits*p.N))
	r.Mod(r, p.P)
	rr := new(big.Int).Mul(r, r)
	rr.Mod(rr, p.P)
	o.line("// %sR is R mod p, 1 in the Montgomery domain.", p.pre)
	o.line("var %sR = %s", p.pre, literal(p, r))
	o.line("// %sRR is R² mod p, which takes an element into the domain.", p.pre)
	o.line("var %sRR = %s", p.pre, literal(p, rr))
	return o.Bytes()
}

// literal is x as one of p's elements, out of the domain.
func literal(p *prime, x *big.Int) string {
	var l []string
	for _, v := range limbsOf(x, p.N) {
		l = append(l, fmt.Sprint(v))
	}
	return p.Type + "{" + strings.Join(l, ", ") + "}"
}
