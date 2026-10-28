package main

import "fmt"

// mul writes p's Montgomery multiplication, scanning the product a column
// at a time.
func mul(o *out, p *prime) {
	n := p.N
	o.line("// %s sets z = x*y/R mod p, by Montgomery multiplication scanning", p.fn("Mul"))
	o.line("// the product a column at a time, each column's sum exact.")
	o.line("func %s(z, x, y *%s) {", p.fn("Mul"), p.Type)
	for i := range n {
		o.line("x%d, y%d := float64(x.l%d), float64(y.l%d)", i, i, i, i)
	}
	o.line("var acc, c, v float64")
	for i := range n {
		t := append(column(p, i, 0), fmt.Sprintf("x%d*y0", i))
		if i > 0 {
			t = append(t, "c")
		}
		o.line("acc = %s", tree(t))
		o.line("v = acc - (acc - half + round - round)")
		if p.PInv != 1 {
			o.line("v *= %d", p.PInv)
			o.line("v -= v - half + round - round")
		}
		o.line("m%d := v", i)
		o.line("c = (acc + %s) * inv", times(p.Limbs[0], fmt.Sprintf("m%d", i)))
	}
	for i := n; i < 2*n; i++ {
		o.line("acc = %s", tree(append(column(p, i, i-n+1), "c")))
		o.line("v = acc - half + round - round")
		o.line("t%d := int32(acc - v)", i-n)
		o.line("c = v * inv")
	}
	reduce(o, p, "int32(c)")
	o.line("}\n")
}

// column is the terms of column i from x[from]: x[j]·y[i-j], and the
// reduction's p[i-j]·m[j], up to j = min(i, n) - 1.
func column(p *prime, i, from int) []string {
	var t []string
	for j := from; j < min(i, p.N); j++ {
		t = append(t, fmt.Sprintf("x%d*y%d", j, i-j))
		if m := times(p.Limbs[i-j], fmt.Sprintf("m%d", j)); m != "" {
			t = append(t, m)
		}
	}
	return t
}

// times is the constant limb l times v, or "" when l is 0.
func times(l int64, v string) string {
	switch l {
	case 0:
		return ""
	case 1:
		return v
	}
	return fmt.Sprintf("%d*%s", l, v)
}
