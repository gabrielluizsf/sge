package dialogue

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func blank(w, h int) *canvas.Canvas {
	cv, err := canvas.NewCanvas(max(w, 1), max(h, 1))
	if err != nil {
		panic(err)
	}
	return cv
}

// flat is an Art of blank pictures, a letter six pixels wide.
type flat struct{}

func (flat) Panel(w, h int, _ bool) *canvas.Canvas { return blank(w*2, h*2) }
func (flat) Icon(string) *canvas.Canvas            { return blank(32, 32) }
func (flat) Text(s string, _ rune) *canvas.Canvas  { return blank(len(s)*6, 12) }

// said is a box saying text with a few buttons, painted once on a 1280 by
// 720 screen.
func said(text string) (*Box, *canvas.Canvas) {
	b := &Box{Speaker: "Ana", Scale: 2, Speed: 30, Choices: []string{"Conversar", "Dar cenoura", "Tchau"}}
	b.Say(text)
	cv := blank(1280, 720)
	b.Paint(cv, flat{}, blank(80, 100))
	return b, cv
}

// The box sits centred at the bottom of the screen, short.
func TestTheBoxIsAtTheBottom(t *testing.T) {
	b, cv := said("Oi!")
	f := b.Frame()
	if f.Y+f.H > cv.Height || f.Y < cv.Height/2 || f.X+f.W/2 != cv.Width/2 {
		t.Errorf("frame %v on %dx%d", f, cv.Width, cv.Height)
	}
}

// Words appear a letter at a time; a key shows them whole, and the buttons
// can not be chosen until they are said.
func TestWordsAreTypedAndAKeyShowsThemWhole(t *testing.T) {
	b, _ := said("Bom dia, tenho trabalho para você hoje.")
	if b.Done() || b.Ready() {
		t.Fatal("done before it began")
	}
	b.Update(0.1)
	if b.Done() {
		t.Fatal("done after three letters")
	}
	if b.Handle(Input{Select: true}) != Skipped || !b.Done() {
		t.Fatal("a key did not show the page whole")
	}
	if b.Handle(Input{Select: true}) != Chosen {
		t.Error("another did not choose")
	}
}

// Long words fill pages; a key goes on to the next once one is shown, and
// the buttons wait for the last.
func TestLongWordsAreInPages(t *testing.T) {
	b, cv := said(strings.Repeat("palavra e mais palavra ", 40))
	if b.Pages() < 2 {
		t.Fatalf("%d pages", b.Pages())
	}
	b.Handle(Input{Select: true})
	if b.Ready() || b.Handle(Input{Select: true}) != Advanced || b.Page() != 1 {
		t.Fatalf("page %d", b.Page())
	}
	for !b.Last() {
		b.Handle(Input{Select: true})
		b.Handle(Input{Select: true})
	}
	b.Handle(Input{Select: true})
	b.Paint(cv, flat{}, nil)
	if !b.Ready() {
		t.Error("not ready on the last page, shown")
	}
}

// The arrows and the mouse choose among the buttons; a press on one, once
// the words are said, answers.
func TestButtonsAreChosenByKeysAndMouse(t *testing.T) {
	b, _ := said("Oi.")
	b.Update(10)
	b.Handle(Input{Dy: 1})
	if b.Choice != 1 {
		t.Fatalf("choice %d", b.Choice)
	}
	r := b.buttons[2]
	at := Point{X: r.X + 5, Y: r.Y + 5}
	b.Handle(Input{Mouse: at})
	if b.Choice != 2 {
		t.Fatalf("choice %d after the mouse moved over the third", b.Choice)
	}
	b.Choice = 0
	if b.Handle(Input{Mouse: at, Press: true}) != Chosen || b.Choice != 2 {
		t.Errorf("choice %d after a press", b.Choice)
	}
	if b.Handle(Input{Cancel: true}) != Left {
		t.Error("escape did not leave")
	}
}
