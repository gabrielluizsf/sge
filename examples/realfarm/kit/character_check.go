package kit

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/textfield"
)

// MaxName is the most characters a name has.
const MaxName = 20

// maxWear is the most proxies a character may wear.
const maxWear = 16

var (
	knownOnce sync.Once
	known     map[string]bool
)

// inCatalogue reports whether name is a proxy the character creator may
// offer.
func inCatalogue(name string) bool {
	knownOnce.Do(func() {
		known = map[string]bool{}
		for _, n := range human.Catalogue() {
			known[n] = true
		}
	})
	return known[name]
}

// unit reports whether v is a number from 0 to 1.
func unit(v float64) bool { return !math.IsNaN(v) && v >= 0 && v <= 1 }

// colourOK reports whether c is a colour: its channels from 0 to 1, and
// nothing in them when it is not set.
func colourOK(c human.Colour) bool {
	if !unit(c.R) || !unit(c.G) || !unit(c.B) {
		return false
	}
	return c.Set || (c.R == 0 && c.G == 0 && c.B == 0)
}

// Check says why c is not a character the creator could make, nil if it
// is: every number in range, every colour a colour, every file one of the
// catalogue's. A file read from disk is checked before it is used.
func (c Character) Check() error {
	b := c.Body
	for _, v := range []float64{b.Gender, b.Age, b.African, b.Asian, b.European, b.Muscle, b.Weight, b.Height, b.Proportions} {
		if !unit(v) {
			return fmt.Errorf("character: body number %v out of range", v)
		}
	}
	if c.Name != "" && !textfield.Valid(c.Name, MaxName) {
		return fmt.Errorf("character: a name of %d bytes that is not one", len(c.Name))
	}
	l := c.Look
	if l.Skin != "" {
		return errors.New("character: a skin file")
	}
	if l.Eyes != "" && (!strings.HasPrefix(l.Eyes, "eyes/materials/") || !inCatalogue(l.Eyes)) {
		return fmt.Errorf("character: unknown eyes %q", l.Eyes)
	}
	if len(l.Wear) > maxWear || len(l.ClothesColours) > len(l.Wear) {
		return errors.New("character: too many things worn")
	}
	seen := map[string]bool{}
	for _, w := range l.Wear {
		if !inCatalogue(w) || seen[w] {
			return fmt.Errorf("character: unknown or repeated %q", w)
		}
		seen[w] = true
	}
	if !unit(l.SkinTone) || !colourOK(l.HairColour) || !colourOK(l.EyeColour) {
		return errors.New("character: skin tone or colour out of range")
	}
	for _, cc := range l.ClothesColours {
		if !colourOK(cc) {
			return errors.New("character: clothes colour out of range")
		}
	}
	return nil
}
