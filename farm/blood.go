package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/particle"
	"github.com/gabrielluizsf/sge/vector"
)

// bloodSpecks is how many specks of blood the farm keeps at most.
const bloodSpecks = 256

// bloodKind is what blood is to the particle system: it spurts up and out,
// falls, and lies on the ground a while before it dries and fades.
var bloodKind = particle.Config{
	Speed: 40 * scale, Lift: 90 * scale, Gravity: 420 * scale,
	Rest: 6, Spread: 2, Size: scale, Grow: 2 * scale,
	Colors: []canvas.Color{0xFFB3141C, 0xFF8E0F18, 0xFFC8222A, 0xFF6E0A12},
}

// newBlood is the two objects that show the farm's blood: the drops in the
// air over everything, the stains on the ground under what stands on it.
func (f *farm) newBlood() (air, stains *game.Object) {
	f.blood = particle.New(bloodSpecks, bloodKind, 7)
	air = game.NewObject("blood")
	air.Z = zAir
	air.AddScript(f.blood)
	stains = game.NewObject("stains")
	stains.Z = zStains
	stains.AddScript(particle.Lying{System: f.blood})
	f.bloodObj, f.stainObj = air, stains
	return air, stains
}

// bleed throws n specks of blood from the body at, as high as the chest,
// away from the blow that came from from.
func (f *farm) bleed(at, from vector.Vector3, n int) {
	away := at.Sub(from)
	if l := away.Length(); l > 0 {
		away = away.Scale(40 * scale / l)
	}
	f.blood.Emit(vector.Vector3{X: at.X, Y: at.Y, Z: 0.7 * tile}, away, n)
}
