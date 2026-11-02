package backpack

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// bagOf is a Source of a few items: seeds, more seeds, a hoe and a crop,
// which is carried but can not be held.
type bagOf struct{}

var things = []Item{
	{ID: "a_seeds", Name: "Sementes A", Icon: "a", Count: 5, Hand: true},
	{ID: "b_seeds", Name: "Sementes B", Icon: "b", Count: 2, Hand: true},
	{ID: "hoe", Name: "Enxada", Icon: "hoe", Count: -1, Info: []string{"Limpa o mato"}, Hand: true},
	{ID: "corn", Name: "Milho", Icon: "corn", Count: 3, Info: []string{"Vende por $9"}},
}

func (bagOf) Items() []Item { return things }
func (bagOf) Find(id string) (Item, bool) {
	for _, it := range things {
		if it.ID == id {
			return it, true
		}
	}
	return Item{}, false
}

// flat is an Art that draws a blank picture of any size.
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

// opened is a pack of four slots on a screen 1280 by 720, open, with the
// hoe on slot 1.
func opened() *Pack {
	p := New(bagOf{}, 4, 8)
	p.Scale = 2
	p.Hotbar[0] = "hoe"
	p.Place(View{W: 1280, H: 720, Scale: 2})
	p.Show()
	return p
}

// centre is the middle of slot s.
func centre(p *Pack, s Spot) Point {
	r := p.rectOf(s)
	return Point{r.X + r.W/2, r.Y + r.H/2}
}

// Arrows move the cursor over the grid and down to the hotbar, and back.
func TestTheArrowsMoveTheCursorFromGridToBar(t *testing.T) {
	p := opened()
	p.Handle(Input{Dx: 1})
	if p.Cursor() != (Spot{Grid, 1}) {
		t.Fatalf("cursor %v after right", p.Cursor())
	}
	p.Handle(Input{Dx: -1})
	p.Handle(Input{Dx: -1}) // the wall
	if p.Cursor() != (Spot{Grid, 0}) {
		t.Fatalf("cursor %v at the wall", p.Cursor())
	}
	for range 3 {
		p.Handle(Input{Dy: 1})
	}
	if p.Cursor() != (Spot{Bar, 0}) {
		t.Fatalf("cursor %v after going down the grid", p.Cursor())
	}
	p.Handle(Input{Dx: 1})
	p.Handle(Input{Dy: -1})
	if c := p.Cursor(); c.Area != Grid || c.Index != 2*8+1 {
		t.Errorf("cursor %v after going up", c)
	}
}

// Enter takes up the item under the cursor, and puts it down on a hotbar
// slot, which then holds it; the slot's old item is not lost, it trades
// if the new one was on the bar.
func TestEnterPicksAnItemAndPutsItOnASlot(t *testing.T) {
	p := opened()
	if ev := p.Handle(Input{Select: true}); ev != Picked || p.Held() != "a_seeds" {
		t.Fatalf("event %v, held %q", ev, p.Held())
	}
	for range 3 {
		p.Handle(Input{Dy: 1})
	}
	if ev := p.Handle(Input{Select: true}); ev != Bound {
		t.Fatalf("event %v putting it down", ev)
	}
	if p.Hotbar[0] != "a_seeds" || p.Held() != "" {
		t.Errorf("hotbar %v, held %q", p.Hotbar, p.Held())
	}
}

// A number key puts the item under the cursor on that slot, and the slot's
// previous item leaves it; a crop can not be held, so it is refused.
func TestANumberKeyAssignsToThatSlot(t *testing.T) {
	p := opened()
	p.Handle(Input{Dx: 1})
	if ev := p.Handle(Input{Number: 1}); ev != Bound || p.Hotbar[0] != "b_seeds" {
		t.Fatalf("event %v, hotbar %v", ev, p.Hotbar)
	}
	p.Handle(Input{Dx: 2}) // the corn
	if ev := p.Handle(Input{Number: 2}); ev != Refused || p.Hotbar[1] != "" {
		t.Errorf("event %v, hotbar %v for a crop", ev, p.Hotbar)
	}
}

// An item is on one slot at most: put on another, the slots trade.
func TestAnItemIsOnOneSlotAtMost(t *testing.T) {
	p := opened()
	p.Bind("a_seeds", 2)
	p.Bind("hoe", 2)
	if p.Hotbar[2] != "hoe" || p.Hotbar[0] != "a_seeds" || p.Slot("hoe") != 2 {
		t.Errorf("hotbar %v", p.Hotbar)
	}
}

// Dragging with the mouse from the grid to a slot binds it; dragging a
// slot's item off the bar takes it off.
func TestDraggingWithTheMouse(t *testing.T) {
	p := opened()
	from, to := centre(p, Spot{Grid, 1}), centre(p, Spot{Bar, 2})
	if ev := p.Handle(Input{Mouse: from, Press: true}); ev != Picked {
		t.Fatalf("press: %v", ev)
	}
	p.Handle(Input{Mouse: to})
	if ev := p.Handle(Input{Mouse: to, Release: true}); ev != Bound || p.Hotbar[2] != "b_seeds" {
		t.Fatalf("release: %v, hotbar %v", ev, p.Hotbar)
	}
	off := centre(p, Spot{Bar, 2})
	p.Handle(Input{Mouse: off, Press: true})
	nowhere := Point{5, 5}
	if ev := p.Handle(Input{Mouse: nowhere, Release: true}); ev != Unbound || p.Hotbar[2] != "" {
		t.Errorf("dragged off: %v, hotbar %v", ev, p.Hotbar)
	}
}

// Pressing and letting go on one slot is a click, which leaves the item
// up for the next click to put down.
func TestAClickThenAClick(t *testing.T) {
	p := opened()
	at := centre(p, Spot{Grid, 2})
	p.Handle(Input{Mouse: at, Press: true})
	p.Handle(Input{Mouse: at, Release: true})
	if p.Held() != "hoe" {
		t.Fatalf("held %q after a click", p.Held())
	}
	to := centre(p, Spot{Bar, 3})
	if ev := p.Handle(Input{Mouse: to, Press: true}); ev != Bound || p.Hotbar[3] != "hoe" || p.Hotbar[0] != "" {
		t.Errorf("second click: %v, hotbar %v", ev, p.Hotbar)
	}
}

// Escape lets go of what is held first, and closes the screen next.
func TestEscapeLetsGoThenCloses(t *testing.T) {
	p := opened()
	p.Handle(Input{Select: true})
	if ev := p.Handle(Input{Cancel: true}); ev != Dropped || !p.Open {
		t.Fatalf("first escape: %v, open %v", ev, p.Open)
	}
	if ev := p.Handle(Input{Cancel: true}); ev != Closed || p.Open {
		t.Errorf("second escape: %v, open %v", ev, p.Open)
	}
	if ev := p.Handle(Input{Select: true}); ev != Nothing {
		t.Errorf("a closed pack did %v", ev)
	}
}

// The pack draws, closed and open, held item and tooltip included.
func TestThePackPaints(t *testing.T) {
	p := opened()
	cv := blank(1280, 720)
	p.Paint(cv, flat{})
	p.Handle(Input{Select: true})
	p.Paint(cv, flat{})
	p.Hide()
	p.Paint(cv, flat{})
}
