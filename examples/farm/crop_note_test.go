package farm

import "testing"

// A crop's note tells the days left and whether it is watered.
func TestCropNote(t *testing.T) {
	g, f := newGame(idle{})
	_ = g
	p := sown(f, 5) // corn
	if got := cropNote(p); got != "Milho - brota em 1 dia - precisa de água" {
		t.Fatalf("dry seed: %q", got)
	}
	water(f, p)
	if got := cropNote(p); got != "Milho - brota em 1 dia - regado" {
		t.Fatalf("watered seed: %q", got)
	}
	p.days, p.stage = 4, stageOf(crops[5], 4)
	if got := cropNote(p); got[:17] != "Milho - cresce em" {
		t.Fatalf("growing: %q", got)
	}
	p.stage = ripe
	if got := cropNote(p); got != "Milho - pronto para colher" {
		t.Fatalf("ripe: %q", got)
	}
}

// Only a planted plot near the farmer is found.
func TestNearPlot(t *testing.T) {
	g, f := newGame(idle{})
	p := sown(f, 5)
	o := g.Find("farmer")
	o.Transform.Position = standIn(f, p.col, p.row+1)
	if f.nearPlot(o) != p {
		t.Fatal("the planted plot beside the farmer is not found")
	}
	o.Transform.Position = standIn(f, p.col, p.row+9)
	if f.nearPlot(o) != nil {
		t.Fatal("a plot far away was found")
	}
}
