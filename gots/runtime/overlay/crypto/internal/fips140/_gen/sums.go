package main

import "fmt"

// add writes p's sum, reduced.
func add(o *out, p *prime) {
	o.line("func %s(z, x, y *%s) {", p.fn("Add"), p.Type)
	o.line("var v, c int32")
	for i := range p.N {
		o.line("v = x.l%d + y.l%d%s", i, i, plusIf(i > 0, " + c"))
		o.line("c = v >> limbBits")
		o.line("t%d := v & mask", i)
	}
	reduce(o, p, "c")
	o.line("}\n")
}

// sub writes p's difference: p is added back when it is below zero.
func sub(o *out, p *prime) {
	o.line("func %s(z, x, y *%s) {", p.fn("Sub"), p.Type)
	o.line("var v, b int32")
	for i := range p.N {
		o.line("v = x.l%d - y.l%d%s", i, i, plusIf(i > 0, " - b"))
		o.line("b = -(v >> limbBits)")
		o.line("d%d := v & mask", i)
	}
	o.line("// Below zero: add p back.")
	o.line("k := -b")
	o.line("c := int32(0)")
	for i := range p.N {
		pl := ""
		if l := p.Limbs[i]; l != 0 {
			pl = fmt.Sprintf(" + %d&k", l)
		}
		o.line("v = d%d%s + c", i, pl)
		if i < p.N-1 {
			o.line("c = v >> limbBits")
		}
		o.line("z.l%d = v & mask", i)
	}
	o.line("}\n")
}

func plusIf(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}
