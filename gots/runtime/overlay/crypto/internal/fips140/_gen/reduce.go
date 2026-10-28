package main

// reduce writes z = t mod p, for t, t0 to t(n-1) and the carry top,
// below 2p: t - p, unless that borrows.
func reduce(o *out, p *prime, top string) {
	o.line("var d, b int32")
	for i := range p.N {
		o.line("d = t%d - %d%s", i, p.Limbs[i], plusIf(i > 0, " - b"))
		o.line("b = -(d >> limbBits)")
		o.line("s%d := d & mask", i)
	}
	o.line("// k is all ones when t is below p, and t is kept.")
	o.line("k := (%s - b) >> limbBits", top)
	for i := range p.N {
		o.line("z.l%d = s%d ^ ((s%d ^ t%d) & k)", i, i, i, i)
	}
}
