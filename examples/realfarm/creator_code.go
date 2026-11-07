package main

import (
	"github.com/gabrielluizsf/sge/colourcode"
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/textfield"
)

// colourOfCode is the colour a code spells; alpha has no place in a
// garment's tint, so it is told apart by the caller, not carried here.
func colourOfCode(c colourcode.RGBA) human.Colour {
	return human.RGB(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
}

// codeOf is the code of colour c, "" when it is not set.
func codeOf(c human.Colour) string {
	if !c.Set {
		return ""
	}
	to8 := func(v float64) uint8 { return uint8(min(max(v, 0), 1)*255 + 0.5) }
	return colourcode.Format(colourcode.RGBA{R: to8(c.R), G: to8(c.G), B: to8(c.B), A: 255})
}

// startCode gives the keyboard to the code field of the row chosen, with
// the code of its colour in it.
func (c *creator) startCode() {
	c.typing = true
	c.code = codeOf(*c.at().field(c.st))
	c.note, c.fresh = "", true
}

// typeCode types this frame into the code field and leaves it on Enter
// with a good code, or on Escape or the arrows up and down. A code that
// reads is worn at once; a bad one on Enter changes nothing and says so.
func (c *creator) typeCode(in cin) cin {
	field := c.at().field(c.st)
	if in.Back && c.code != "" {
		c.fresh = false
		c.code = textfield.Backspace(c.code)
		c.note = ""
	}
	if in.Typed != "" {
		if c.fresh {
			c.code, c.fresh = "", false // typing over the shown code replaces it
		}
		c.code = colourcode.Insert(c.code, in.Typed)
		c.note = ""
	}
	if in.Back || in.Typed != "" {
		if rgba, err := colourcode.Parse(c.code); err == nil {
			*field = colourOfCode(rgba)
			if rgba.A != 255 {
				c.note = i18n.T("Clothes cannot be see-through: transparency is ignored")
			}
		}
	}
	if in.Select {
		in.Select = false
		if _, err := colourcode.Parse(c.code); err != nil {
			c.note = i18n.T("Not a colour code. Try #RRGGBB or rgba(r, g, b, a)")
		} else {
			c.typing = false
		}
	}
	if in.Dy != 0 || in.Cancel {
		c.typing, c.note = false, ""
	}
	return in
}
