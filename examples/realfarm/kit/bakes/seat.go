package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// ChairModel is the asset the player sits on.
const ChairModel = "painted_wooden_chair_02"

// seatForward is how far forward of the spot under the hips, in metres, the
// chair's middle stands: with the hips over the back half of its seat, the
// back against its back.
const seatForward = 0.05

// seatOnChair bakes the chair with m, a person already animated, if the
// person sits: the chair at its size in the game, facing the way they do,
// its seat under them. They are one picture, so the chair hides the body
// where it should on every side.
func (b *Baker) seatOnChair(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Sit) {
		return nil
	}
	a, ok := b.chairAsset()
	if !ok {
		return fmt.Errorf("no asset %s to sit on", ChairModel)
	}
	chair, err := Load(a)
	if err != nil {
		return fmt.Errorf("%s: %w", a.ID, err)
	}
	return human.Seat(m, chair, human.SeatPlace{At: [3]float64{0, 0, seatForward}, Scale: a.Scale, During: []string{human.Sit.Name}})
}

// chairAsset is the manifest's chair, and whether it has one.
func (b *Baker) chairAsset() (kit.Asset, bool) {
	for _, a := range b.Manifest.Assets {
		if a.ID == ChairModel {
			return a, true
		}
	}
	return kit.Asset{}, false
}
