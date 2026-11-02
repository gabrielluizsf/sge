package menu

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// flat is an Art that draws blank pictures, a letter six pixels wide.
type flat struct{}

func blank(w, h int) *canvas.Canvas {
	cv, err := canvas.NewCanvas(max(w, 1), max(h, 1))
	if err != nil {
		panic(err)
	}
	return cv
}

func (flat) Panel(w, h int, _ bool) *canvas.Canvas { return blank(w*2, h*2) }
func (flat) Icon(string) *canvas.Canvas            { return blank(32, 32) }
func (flat) Text(s string, _ rune) *canvas.Canvas  { return blank(len(s)*6, 12) }

func pause() *List {
	l := &List{Title: "Pausa", Scale: 2, MinW: 120, Lines: []Line{{Text: "Continuar"}, {Text: "Créditos"}, {Text: "Sair"}}}
	l.Paint(blank(1280, 720), flat{})
	return l
}

// The arrows go round the lines, Enter chooses, Escape leaves.
func TestTheListIsChosenFromByKeys(t *testing.T) {
	l := pause()
	if l.Handle(Input{Dy: -1}) != Moved || l.Choice != 2 {
		t.Fatalf("choice %d after up from the first", l.Choice)
	}
	l.Handle(Input{Dy: 1})
	l.Handle(Input{Dy: 1})
	if l.Handle(Input{Select: true}) != Chosen || l.Choice != 1 {
		t.Errorf("choice %d", l.Choice)
	}
	if l.Handle(Input{Cancel: true}) != Cancelled {
		t.Error("escape did not leave")
	}
}

// The mouse chooses the line it moves over and picks the one it presses.
func TestTheListIsChosenFromByTheMouse(t *testing.T) {
	l := pause()
	r := l.rows[2]
	at := Point{X: r.X + 10, Y: r.Y + 5}
	if l.Handle(Input{Mouse: at}) != Moved || l.Choice != 2 {
		t.Fatalf("choice %d after moving over the third", l.Choice)
	}
	l.Choice = 0
	if l.Handle(Input{Mouse: at, Press: true}) != Chosen || l.Choice != 2 {
		t.Errorf("choice %d after pressing it", l.Choice)
	}
}

// A list of more lines than fit scrolls to keep the choice in sight and by
// the wheel.
func TestALongListScrolls(t *testing.T) {
	l := &List{Title: "Créditos", Scale: 2, MinW: 120}
	for range 60 {
		l.Lines = append(l.Lines, Line{Text: "a", Sub: "b", Button: "ver"})
	}
	cv := blank(1280, 720)
	l.Paint(cv, flat{})
	if l.shown >= 60 || l.shown < 5 {
		t.Fatalf("%d lines in sight", l.shown)
	}
	for range l.shown + 3 {
		l.Handle(Input{Dy: 1})
		l.Paint(cv, flat{})
	}
	if l.Top == 0 || l.Choice < l.Top || l.Choice >= l.Top+l.shown {
		t.Errorf("top %d, choice %d, %d in sight", l.Top, l.Choice, l.shown)
	}
	l.Handle(Input{Wheel: -50})
	if l.Top != 0 || l.Choice >= l.shown {
		t.Errorf("top %d, choice %d after the wheel", l.Top, l.Choice)
	}
}

// The reader wraps a long text, scrolls it and never past its end.
func TestTheReaderScrollsALicence(t *testing.T) {
	r := &Reader{Scale: 2}
	r.Open("Licença", strings.Repeat("palavra outra coisa mais ", 400)+"\n\nfim")
	cv := blank(1280, 720)
	r.Paint(cv, flat{})
	if len(r.wrapped) < 20 {
		t.Fatalf("%d lines", len(r.wrapped))
	}
	r.Handle(Input{Page: 1})
	if r.Top != r.shown-1 {
		t.Errorf("top %d after a page", r.Top)
	}
	r.Handle(Input{Wheel: 100000})
	r.Paint(cv, flat{})
	if r.Top != len(r.wrapped)-r.shown {
		t.Errorf("top %d at the end of %d lines", r.Top, len(r.wrapped))
	}
	if r.Handle(Input{Cancel: true}) != Cancelled {
		t.Error("escape did not close")
	}
}
