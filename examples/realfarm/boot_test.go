package main

import (
	"fmt"
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm"
)

// The loading screen is written in the HUD's serif, not in the stylized
// farm's face it once borrowed.
func TestTheLoadingScreenIsInTheRealisticFace(t *testing.T) {
	got := bootLabel(3, 9, "tree ipe")
	if !same(got, serifText("Preparando o mundo 3/9 - tree ipe", 'w')) {
		t.Error("the loading line is not in Source Serif")
	}
	if same(got, farm.SVGArt().Text(fmt.Sprintf("Preparando o mundo %d/%d - %s", 3, 9, "tree ipe"), 'w')) {
		t.Error("the loading line is in the stylized face")
	}
}
