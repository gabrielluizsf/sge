package human

import (
	"strings"

	"github.com/gabrielluizsf/sge/random"
)

// fabrics are the colours of cloth: muted, and some bright.
var fabrics = []Colour{
	RGB(0.65, 0.18, 0.15), RGB(0.85, 0.55, 0.15), RGB(0.82, 0.72, 0.30), RGB(0.30, 0.50, 0.25),
	RGB(0.15, 0.45, 0.55), RGB(0.18, 0.28, 0.55), RGB(0.45, 0.25, 0.50), RGB(0.80, 0.45, 0.55),
	RGB(0.85, 0.82, 0.75), RGB(0.45, 0.35, 0.25), RGB(0.30, 0.30, 0.32), RGB(0.55, 0.55, 0.50),
}

// colourClothes is a colour for each of wear's tops, bottoms and dresses,
// three in four, by place; nil if none. Suits, shoes and hats keep the
// pack's.
func colourClothes(wear []string, r *random.Source) []Colour {
	out := make([]Colour, len(wear))
	any := false
	for i, name := range wear {
		if !strings.HasPrefix(name, "clothes/") || !colourable(name) {
			continue
		}
		if r.Float() < 0.75 {
			out[i], any = fabrics[r.Intn(len(fabrics))], true
		}
	}
	if !any {
		return nil
	}
	return out
}

// colourable is whether the proxy is a top, bottom or dress.
func colourable(name string) bool {
	for _, list := range [][]garment{tops, bottoms, dresses} {
		for _, g := range list {
			if g.item == name {
				return true
			}
		}
	}
	return false
}
