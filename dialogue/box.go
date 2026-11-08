// Package dialogue is the box a role-playing game has its people talk in:
// a short panel at the bottom of the screen with the speaker's portrait on
// its left, their name over their words, which appear a letter at a time
// and, if long, in pages, and what the player may answer as buttons on its
// right. Everything a character says goes through it.
package dialogue

import "github.com/gabrielluizsf/sge/menu"

// Art is how the game draws the box.
type Art = menu.Art

// Point and Rect are pixels, and parts of the screen.
type (
	Point = menu.Point
	Rect  = menu.Rect
)

// Input is what the player did this frame: arrows as steps along the
// buttons, Enter or space, Escape, and the mouse.
type Input struct {
	Dy     int
	Select bool
	Cancel bool
	Mouse  Point
	Press  bool
}

// Result is what an input did to the box.
type Result int

// The results.
const (
	Nothing  Result = iota
	Moved           // another button was chosen
	Skipped         // the page was shown whole
	Advanced        // the next page began
	Chosen          // a button was chosen: see Box.Choice
	Left            // the player left
)

// Box is a conversation: who speaks, what they say, and what may be
// answered.
type Box struct {
	// Speaker is the name over the words.
	Speaker string
	// Choices are the buttons, Choice the one chosen.
	Choices []string
	Choice  int
	// Scale is how many pixels a pixel of the game's art takes.
	Scale int
	// Speed is how many letters appear a second; 0 shows a page whole.
	Speed float64

	text    string
	pages   [][]string
	page    int
	shown   float64
	laid    bool
	words   map[string]int
	frame   Rect
	buttons []Rect // where the buttons in sight are
	first   int    // the first button in sight
	at      Point
}

// Say begins to say text, from its first letter, in place of what was
// being said.
func (b *Box) Say(text string) {
	b.text, b.pages, b.page, b.shown, b.laid = text, nil, 0, 0, false
}

// Text is what is being said.
func (b *Box) Text() string { return b.text }

// Frame is where the box lies on the screen, as of the last Paint.
func (b *Box) Frame() Rect { return b.frame }

// Pages is how many pages the words fill, as of the last Paint.
func (b *Box) Pages() int { return len(b.pages) }

// Page is the page being said, from 0.
func (b *Box) Page() int { return b.page }

// Update shows dt seconds more of the words.
func (b *Box) Update(dt float64) {
	if b.Speed <= 0 {
		b.shown = 1 << 30
		return
	}
	b.shown += dt * b.Speed
}

// letters is how many letters the page has.
func (b *Box) letters() int {
	n := 0
	if b.page < len(b.pages) {
		for _, l := range b.pages[b.page] {
			n += len([]rune(l))
		}
	}
	return n
}

// Done reports whether the whole page is shown.
func (b *Box) Done() bool { return b.laid && int(b.shown) >= b.letters() }

// Last reports whether the page is the last one.
func (b *Box) Last() bool { return b.page >= len(b.pages)-1 }

// Ready reports whether the words are all said, and the buttons may be
// chosen.
func (b *Box) Ready() bool { return b.Done() && b.Last() }
