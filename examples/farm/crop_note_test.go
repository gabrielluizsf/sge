package farm

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/vector"
)

// A crop's note tells the days left and whether it is watered.
func TestCropNote(t *testing.T) {
	g, f := newGame(idle{})
	_ = g
	p := sown(f, 5) // corn
	if got := cropNote(p); got != "Corn - sprouts in 1 day - needs water" {
		t.Fatalf("dry seed: %q", got)
	}
	water(f, p)
	if got := cropNote(p); got != "Corn - sprouts in 1 day - watered" {
		t.Fatalf("watered seed: %q", got)
	}
	p.days, p.stage = 4, stageOf(crops[5], 4)
	if got := cropNote(p); !strings.HasPrefix(got, "Corn - grows in") {
		t.Fatalf("growing: %q", got)
	}
	p.stage = ripe
	if got := cropNote(p); got != "Corn - ready to harvest" {
		t.Fatalf("ripe: %q", got)
	}
}

// The note is of the plot the cursor marks, the cell before the farmer or
// the one the mouse is on, and of no other: not of one beside them.
func TestTheNoteIsOfTheMarkedPlot(t *testing.T) {
	g, f := newGame(idle{})
	p := sown(f, 5)
	o := g.Find("farmer")
	o.Transform.Position = standIn(f, p.col, p.row+1)
	f.hands.facing = vector.Vector3{Y: -1} // toward the plot
	if f.markedPlot(o) != p {
		t.Fatal("the plot the cursor marks has no note")
	}
	f.hands.facing = vector.Vector3{X: 1} // beside it, the cursor elsewhere
	if f.markedPlot(o) != nil {
		t.Fatal("a plot the cursor does not mark has a note")
	}
	f.hands.aim = &[2]int{p.col, p.row} // the mouse on it
	if f.markedPlot(o) != p {
		t.Fatal("the plot the mouse marks has no note")
	}
}

// In Portuguese the note reads as the game always did: the language
// reaches the crop's name and every word of the line.
func TestCropNoteInPortuguese(t *testing.T) {
	i18n.Set(i18n.Portuguese)
	defer i18n.Set(i18n.English)
	_, f := newGame(idle{})
	p := sown(f, 5) // corn
	if got := cropNote(p); got != "Milho - brota em 1 dia - precisa de água" {
		t.Errorf("note %q", got)
	}
}
