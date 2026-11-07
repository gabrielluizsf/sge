package main

import (
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/textfield"
)

// typeName types what the player typed this frame into the name, if the
// field is focused, and leaves the field on Enter or the arrows up and
// down; it gives back the input with what it took out. A letter typed is
// a letter and never a key that moves: inputOf does not read the keys
// while the field is focused.
func (c *creator) typeName(in cin) cin {
	if !c.typing {
		return in
	}
	if c.at().field != nil {
		return c.typeCode(in)
	}
	if in.Back {
		c.st.name, c.st.named = textfield.Backspace(c.st.name), true
	}
	if in.Typed != "" {
		if !c.st.named {
			c.st.name = "" // typing over the default name replaces it
		}
		c.st.name, c.st.named = textfield.Insert(c.st.name, in.Typed, kit.MaxName), true
	}
	if in.Select {
		c.typing, in.Select = false, false
	}
	if in.Dy != 0 {
		c.typing = false
	}
	return in
}

// adopt makes st what the screen edits, with a new default name.
func (c *creator) adopt(st *charState) {
	st.nameSeed = c.seed * 2654435761
	st.name = st.defaultName()
	c.st = st
}
