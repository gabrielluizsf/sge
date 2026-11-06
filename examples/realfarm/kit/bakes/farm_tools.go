package bakes

import (
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// holdHoe puts the hoe in the hands of m, a person already animated, if
// the person hoes.
func (b *Baker) holdHoe(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Hoe) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	return human.Hold(m, building.Hoe.Model(palette), human.HoeGrip)
}

// holdPickaxe puts the pickaxe in the hands of m, a person already
// animated, if the person mines.
func (b *Baker) holdPickaxe(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Mine) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	return human.Hold(m, building.Pickaxe.Model(palette), human.MineGrip)
}

// holdCan puts the watering can in the right hand of m, a person already
// animated, and the stream it pours, if the person waters.
func (b *Baker) holdCan(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Water) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	if err := human.Hold(m, building.Can.Model(palette), human.CanGrip); err != nil {
		return err
	}
	return human.Pour(m, building.Stream.Model(palette), human.CanSpill)
}

// hasMotion is whether motions has the motion mo.
func hasMotion(motions []human.Motion, mo human.Motion) bool {
	for _, m := range motions {
		if m.Name == mo.Name {
			return true
		}
	}
	return false
}
