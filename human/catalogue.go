package human

import "fmt"

// brows is how many eyebrows the pack has, and lashes how many lashes.
const (
	brows  = 12
	lashes = 4
)

// eyebrows is the proxy of eyebrow n, from 1.
func eyebrows(n int) string { return fmt.Sprintf("eyebrows/eyebrow%03d/eyebrow%03d.mhclo", n, n) }

// eyelashes is the proxy of lashes n, from 1.
func eyelashes(n int) string { return fmt.Sprintf("eyelashes/eyelashes%02d/eyelashes%02d.mhclo", n, n) }

// eyesProxy is the eyes everyone has.
const eyesProxy = "eyes/high-poly/high-poly.mhclo"

// Catalogue is every file Random may choose, by name in the kit: for
// checking they are in the packs, and CC0.
func Catalogue() []string {
	out := []string{eyesProxy}
	for i := 1; i <= brows; i++ {
		out = append(out, eyebrows(i))
	}
	for i := 1; i <= lashes; i++ {
		out = append(out, eyelashes(i))
	}
	for _, e := range eyeColours {
		out = append(out, e.item)
	}
	for _, g := range wardrobe() {
		out = append(out, g.item)
	}
	return out
}
