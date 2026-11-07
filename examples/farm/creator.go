package main

import (
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/screenpick"
)

// creator is the character creation screen of a new game: a list of
// categories on the left, the person in the middle, seen from the front,
// live, and the rows of the category chosen on the right. The person is
// made in the background as it changes, and turns only when the player
// drags it.
type creator struct {
	art  farm.Art
	cats []category
	st   *charState

	cat, row   int // the row chosen, in category cat
	seed       uint32
	live       *livePreview // nil with nothing to make people with
	stage      *liveStage   // draws the person live
	view       liveView
	asked      string                            // the tag of the character last asked of live
	dt         float64                           // the frame's seconds
	typing     bool                              // a text field has the keyboard
	code       string                            // what the code field has typed
	fresh      bool                              // the code shown is replaced by the first key typed
	note       string                            // a short message under the rows
	drop       *dropper                          // the eyedropper while it takes a colour
	openScreen func() (screenpick.Source, error) // nil: the real screen
	turning    bool                              // the mouse holds the person
	turnX      int
	stageAt    rect // where the person is drawn, as last painted
	shownErr   error

	hits []hit
	drag *hit // the slider track the mouse holds

	done     func(any)
	release  func() // lets the background bakes go on again
	obj      *game.Object
	finished bool
}

// newCreator is a creation screen over a first random person; done gets
// the character chosen. live may be nil.
func newCreator(art farm.Art, live *livePreview, done func(any)) *creator {
	c := &creator{art: art, cats: categories(), live: live, done: done, seed: 1}
	if live != nil {
		c.stage = newLiveStage(bestDevice)
	}
	c.adopt(stateOf(human.Random(human.Spec{Seed: c.seed})))
	c.cat, c.row = 0, 0
	c.refresh()
	return c
}

// randomise makes a new person by chance, whom the rows then show.
func (c *creator) randomise() {
	c.seed++
	c.adopt(stateOf(human.Random(human.Spec{Seed: c.seed * 2654435761})))
}

// finish ends the screen with the character made.
func (c *creator) finish() {
	if c.finished {
		return
	}
	c.finished = true
	c.stopPick()
	if c.live != nil {
		c.live.Close()
	}
	if c.stage != nil {
		c.stage.Close()
	}
	if c.release != nil {
		c.release()
	}
	if c.obj != nil {
		c.obj.Destroy()
	}
	c.done(c.st.character())
}

// refresh asks for the person to be shown if the character changed, and
// has the camera frame what the category chosen changes.
func (c *creator) refresh() {
	c.view.Frame(framesHead(c.cats[c.cat].name))
	ch := c.st.character()
	if tag := ch.Tag(); tag != c.asked && c.live != nil {
		c.asked = tag
		c.live.Ask(ch)
	}
}
