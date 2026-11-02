package human

// Part is a kind of thing a person wears, for choosing among the ones
// that suit them.
type Part int

// The parts: hair, tops, bottoms, dresses, suits (worn alone), shoes and
// hats.
const (
	Hair Part = iota
	Top
	Bottom
	Dress
	Suit
	Shoes
	Hat
)

// partLists are the garments of each part, in the order of Part.
func partList(p Part) []garment {
	return [][]garment{hairStyles, tops, bottoms, dresses, suits, shoes, hats}[p]
}

// Options are the proxies of part p that suit a person of sex, years old
// and of african ancestry (0 to 1): what a character creator may offer.
func Options(p Part, sex Sex, years, african float64) []string {
	var out []string
	for _, g := range suited(partList(p), sex, years, african) {
		out = append(out, g.item)
	}
	return out
}

// PartOf is the part the proxy item is, false for any that is not hair,
// clothes, shoes or a hat.
func PartOf(item string) (Part, bool) {
	for p := Hair; p <= Hat; p++ {
		for _, g := range partList(p) {
			if g.item == item {
				return p, true
			}
		}
	}
	return 0, false
}

// Colourable is whether the proxy item is a top, bottom or dress, whose
// colour a Look may choose.
func Colourable(item string) bool { return colourable(item) }

// Eyebrows are the proxies of the pack's eyebrows.
func Eyebrows() []string {
	out := make([]string, brows)
	for i := range out {
		out[i] = eyebrows(i + 1)
	}
	return out
}

// Eyelashes are the proxies of the lashes a person of sex may have: a
// woman's any of the pack's, a man's the first.
func Eyelashes(sex Sex) []string {
	n := lashes
	if sex == Man {
		n = 1
	}
	out := make([]string, n)
	for i := range out {
		out[i] = eyelashes(i + 1)
	}
	return out
}

// EyesProxy is the eyes everyone has.
func EyesProxy() string { return eyesProxy }

// EyeMaterials are the .mhmat of each colour of eyes.
func EyeMaterials() []string {
	out := make([]string, len(eyeColours))
	for i, e := range eyeColours {
		out[i] = e.item
	}
	return out
}

// HairColours are the colours of hair a person may have, black to blond.
func HairColours() []Colour {
	out := make([]Colour, len(hairShades))
	for i, s := range hairShades {
		out[i] = s.colour
	}
	return out
}

// ClothColours are the colours of cloth.
func ClothColours() []Colour { return append([]Colour(nil), fabrics...) }

// SkinToneOf is the middle skin tone, 0 lightest to 1 darkest, of
// ancestry a.
func SkinToneOf(a Ancestry) float64 {
	sum := a.African + a.Asian + a.European
	if sum <= 0 {
		return 0.4
	}
	return (a.African*0.82 + a.Asian*0.3 + a.European*0.14) / sum
}
