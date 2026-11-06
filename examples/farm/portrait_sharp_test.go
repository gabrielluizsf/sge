package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/life"
)

// A portrait is drawn from the art at its own size, not a small picture
// enlarged: enlarged by whole pixels, every block of bustZoom pixels would
// be one colour, and its edges would look blurred and blocky.
func TestAPortraitIsDrawnAtItsSize(t *testing.T) {
	a := SVGArt()
	who := Who{Woman: true, Role: roleFarmer, Stage: life.Adult, Genes: [5]int{3, 1, 4, 1, 5}}
	cv := a.Portrait(who)
	uneven, seen := 0, 0
	for y := 0; y+bustZoom <= cv.Height; y += bustZoom {
		for x := 0; x+bustZoom <= cv.Width; x += bustZoom {
			if cv.At(x, y).A() == 0 && cv.At(x+bustZoom-1, y+bustZoom-1).A() == 0 {
				continue
			}
			seen++
			first := cv.At(x, y)
			flat := true
			for dy := range bustZoom {
				for dx := range bustZoom {
					flat = flat && cv.At(x+dx, y+dy) == first
				}
			}
			if !flat {
				uneven++
			}
		}
	}
	if seen == 0 || uneven*10 < seen {
		t.Errorf("%d of %d blocks show detail: the portrait is a small picture enlarged", uneven, seen)
	}
}
