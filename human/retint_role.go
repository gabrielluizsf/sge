package human

// role is what colours a material: the skin's tone, the hair's colour
// (hair and brows), the eyes' or, for clothes, the colour at place i of
// the look's ClothesColours.
type role struct {
	kind  roleKind
	place int
}

type roleKind uint8

const (
	roleNone roleKind = iota
	roleSkin
	roleHair
	roleEyes
	roleClothes
)

// roleOf is the role of what the look wears at i, the Wear name.
func roleOf(name string, i int) role {
	switch {
	case hasFolder(name, "hair") || hasFolder(name, "eyebrows"):
		return role{kind: roleHair}
	case hasFolder(name, "eyes"):
		return role{kind: roleEyes}
	case hasFolder(name, "clothes"):
		return role{kind: roleClothes, place: i}
	}
	return role{}
}

// tintFor is the tint look l asks of a material in role r, false if it
// leaves it as it is.
func tintFor(r role, l Look) (tint, bool) {
	switch r.kind {
	case roleSkin:
		if l.SkinTone > 0 {
			return skinTint(l.SkinTone), true
		}
	case roleHair:
		return tintTo(l.HairColour, nil)
	case roleEyes:
		return tintTo(l.EyeColour, irises)
	case roleClothes:
		if r.place < len(l.ClothesColours) {
			t, ok := tintTo(l.ClothesColours[r.place], nil)
			t.flat = clothContrast
			return t, ok
		}
	}
	return tint{}, false
}

// clothContrast is the most the light of a cloth's texture may vary about
// its average when it is recoloured, as a share of it: a washed denim or a
// print is eased to this, so one colour does not look painted in patches.
const clothContrast = 0.3
