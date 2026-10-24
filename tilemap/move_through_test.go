package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// Moving through layers stops at the first solid tile of any of them:
// here a wall on one layer to the left, and one on another to the right.
func TestMoveStopsAtTheSolidTilesOfEveryLayer(t *testing.T) {
	ground, _ := New("ground", tiles(t), 10, 1)
	things, _ := New("things", tiles(t), 10, 1)
	ground.Fill(0, 0, 10, 1, 0)
	ground.Set(1, 0, 1) // a wall on the ground, x 4..8
	things.Set(7, 0, 1) // and one among the things, x 28..32
	o := game.NewObject("walker")
	o.Collider = collider.Box{X: -1, Y: -1, Width: 2, Height: 2}
	o.Transform.Position = vector.Vector3{X: 18, Y: 2}
	if hitX, _ := Move(o, 100, 0, ground, things); !hitX || o.Transform.Position.X > 27.02 || o.Transform.Position.X < 26.98 {
		t.Errorf("walking right stopped at %v (hit %v), not at the things' wall, 27", o.Transform.Position.X, hitX)
	}
	if hitX, _ := Move(o, -100, 0, ground, things); !hitX || o.Transform.Position.X < 8.98 || o.Transform.Position.X > 9.02 {
		t.Errorf("walking left stopped at %v (hit %v), not at the ground's wall, 9", o.Transform.Position.X, hitX)
	}
}
