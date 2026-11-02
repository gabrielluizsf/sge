package bakes

import (
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// holdForge puts the hammer in the right hand of m, a person already
// animated, and the tongs with their iron in the left, if the person
// forges.
func (b *Baker) holdForge(m *gltf.Model, motions []human.Motion) error {
	for _, mo := range motions {
		if mo.Name != human.Forge.Name {
			continue
		}
		palette, err := b.Palette()
		if err != nil {
			return err
		}
		if err := human.Hold(m, building.Hammer.Model(palette), human.HammerGrip); err != nil {
			return err
		}
		return human.Hold(m, building.TongsIron.Model(palette), human.TongsGrip)
	}
	return nil
}
