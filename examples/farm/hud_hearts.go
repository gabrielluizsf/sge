package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/bond"
)

// A heart is fondness of perHeart, up to hearts of them for bond.Most;
// they are drawn heartStep pixels of art apart, heartsW wide in all, and
// heartsAt from the panel's left.
const (
	hearts    = 10
	perHeart  = bond.Most / hearts
	heartStep = 12
	heartsW   = heartStep*(hearts-1) + 16
	heartsAt  = 84
)

// bondLabels are how far a bond has come, as the talk panel says it.
var bondLabels = []string{"Desconhecidos", "Conhecidos", "Amigos", "Namorando", "Noivos", "Casados"}

// heartsLabel is how far the player's bond with v has come, in words.
func heartsLabel(v *villager) string { return bondLabels[v.bond.Stage()] }

// hearts draws how fond v is of the player, a full heart for each
// perHeart of fondness and empty ones to make up ten, on the row of the
// talk panel at.
func (h *hud) hearts(cv *canvas.Canvas, v *villager, at [2]int) {
	full := v.bond.Affection / perHeart
	for i := range hearts {
		name := "icon_heart_empty"
		if i < full {
			name = "icon_heart"
		}
		cv.Blit(at[0]+(heartsAt+i*heartStep)*scale, at[1], h.icon(name))
	}
}
