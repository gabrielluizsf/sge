package human

// Look is how a person looks beyond their body: Skin, a .mhmat, or ""
// for MakeHuman's skin of their age, sex and ancestry, blended; Eyes, a
// .mhmat for their eyes' colour, or "" for the eyes' own; and Wear, the
// .mhclo proxies they wear: eyes, brows, lashes, hair, clothes and shoes.
// Names are of files in the kit's FS.
type Look struct {
	Skin, Eyes string
	Wear       []string
}
