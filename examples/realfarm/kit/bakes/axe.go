package bakes

import (
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// holdAxe puts the axe in the hands of m, a person already animated, if
// the person chops.
func (b *Baker) holdAxe(m *gltf.Model, motions []human.Motion) error {
	for _, mo := range motions {
		if mo.Name != human.Chop.Name {
			continue
		}
		palette, err := b.Palette()
		if err != nil {
			return err
		}
		return human.Hold(m, building.AxeIron.Model(palette), human.ChopGrip)
	}
	return nil
}
