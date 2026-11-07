package main

import (
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/screenpick"
)

// dropper is the eyedropper while it is taking a colour: the screen it
// reads, the colour it is for, and the colour under the mouse now.
type dropper struct {
	src    screenpick.Source
	field  *human.Colour
	sawUp  bool // the button has been seen up, so the click that began it is over
	at     screenpick.Reading
	loaded bool // at has been read once
}

// startPick begins the eyedropper for the colour of the row chosen; if
// the screen cannot be read here it says so and nothing else happens.
func (c *creator) startPick() {
	c.typing = false
	open := c.openScreen
	if open == nil {
		open = screenpick.Open
	}
	src, err := open()
	if err != nil {
		c.note = i18n.T("The eyedropper is not available on this system")
		return
	}
	c.note = ""
	c.drop = &dropper{src: src, field: c.at().field(c.st)}
}

// stopPick ends the eyedropper and lets go of the screen.
func (c *creator) stopPick() {
	if c.drop != nil {
		c.drop.src.Close()
		c.drop = nil
	}
}

// pickStep is a frame of the eyedropper: Escape cancels it; the first
// press of the left button after it began, anywhere on the screen, takes
// the colour under the mouse. Nothing is worn until then.
func (c *creator) pickStep(in cin) {
	d := c.drop
	if in.Cancel {
		c.stopPick()
		return
	}
	r, err := d.src.Read()
	if err != nil {
		c.note = i18n.T("The eyedropper could not read the screen")
		c.stopPick()
		return
	}
	d.at, d.loaded = r, true
	switch {
	case !r.Down:
		d.sawUp = true
	case d.sawUp:
		*d.field = human.RGB(float64(r.R)/255, float64(r.G)/255, float64(r.B)/255)
		c.stopPick()
	}
}
