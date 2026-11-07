package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// faceTo is the way, of the four a controller faces, from at toward who;
// none when there is no one.
func faceTo(at vector.Vector3, who *game.Object) string {
	if who == nil {
		return ""
	}
	d := feet(who).Sub(at)
	if math.Abs(d.X) > math.Abs(d.Y) {
		return map[bool]string{true: "right", false: "left"}[d.X > 0]
	}
	return map[bool]string{true: "down", false: "up"}[d.Y > 0]
}
