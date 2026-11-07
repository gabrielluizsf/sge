package farm

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/pathfind"
)

// walkEast is a game with the farmer set down on the path, west of the
// town to the east, run east frame by frame.
func walkEast(t *testing.T, patient bool) func() []canvas.Color {
	g, f := newGame(idle{})
	if !patient { // every picture made at once, none held back
		g.LazyReach = 0
		f.warm.PerFrame = 1 << 20
	}
	cv, err := canvas.NewCanvas(640, 360)
	if err != nil {
		t.Fatal(err)
	}
	o := g.Find("farmer")
	o.Transform.Position = f.center(pathfind.Cell{Col: 34, Row: pathRow})
	g.Draw(cv)
	return func() []canvas.Color {
		o.Transform.Position.X += runSpeed / 60
		g.Step(frame)
		g.Draw(cv)
		return slices.Clone(cv.Pixels)
	}
}

// Pictures made ahead a few a frame, and sprites out of the view left
// unmade, show exactly the frames that making every picture at once does,
// all the way through a town.
func TestFramesAreTheSameWithPicturesMadeAheadOrAtOnce(t *testing.T) {
	patient, hasty := walkEast(t, true), walkEast(t, false)
	for i := range 330 {
		a, b := patient(), hasty()
		if i%10 == 0 && !slices.Equal(a, b) {
			t.Fatalf("frame %d differs", i)
		}
	}
}
