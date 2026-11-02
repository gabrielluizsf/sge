package human

// Look is how a person looks beyond their body: Skin, a .mhmat, or ""
// for MakeHuman's skin of their age, sex and ancestry, blended; Eyes, a
// .mhmat for their eyes' colour, or "" for the eyes' own; and Wear, the
// .mhclo proxies they wear: eyes, brows, lashes, hair, clothes and shoes.
// Names are of files in the kit's FS.
//
// The rest recolours what the pack gives, its texture's detail kept; at
// their zero values nothing is changed. SkinTone is from above 0, the
// lightest skin, to 1, the darkest, blending any skin of the pack's.
// HairColour is the colour of the hair and brows, before age greys them;
// EyeColour of the irises; ClothesColours the colour of each item of
// Wear, by place, the clothes only: unset where it is short.
type Look struct {
	Skin, Eyes     string
	Wear           []string
	SkinTone       float64
	HairColour     Colour
	EyeColour      Colour
	ClothesColours []Colour
}
