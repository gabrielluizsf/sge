package bakes

import (
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// The facings a person milks in, of the baker's eight: to the screen's
// right and left, the cow in profile before them. No other is ever shown,
// and each is a dear thing to bake for the player one makes.
const (
	MilkFacingRight = 6
	MilkFacingLeft  = 2
)

// facesIn reports whether motion mo is baked facing d.
func facesIn(mo human.Motion, d int) bool {
	if mo.Name == human.Milk.Name || mo.Name == human.MilkDown.Name {
		return d == MilkFacingRight || d == MilkFacingLeft
	}
	return true
}

// milkOnStool gives m, a person already animated, the stool under them and
// the pail in front while they milk, if they do.
func (b *Baker) milkOnStool(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Milk) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	return human.Seat(m, building.MilkStool.Model(palette), human.SeatPlace{Scale: 1, During: []string{human.MilkDown.Name, human.Milk.Name}})
}
