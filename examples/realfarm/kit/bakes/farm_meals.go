package bakes

import (
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// holdMorsel puts the bite of food in the right hand of m, a person
// already animated, if the person eats.
func (b *Baker) holdMorsel(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Eat) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	return human.Hold(m, building.Morsel.Model(palette), human.MorselGrip)
}

// holdFillingCan puts the watering can in the right hand of m, a person
// already animated, for the well and the river, and the water that runs
// into it at the well, if the person fills it.
func (b *Baker) holdFillingCan(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Well) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	if err := human.Hold(m, building.Can.Model(palette), human.FillGrip); err != nil {
		return err
	}
	return human.Fill(m, building.Stream.Model(palette), human.WellFall)
}

// holdBottle puts the milk bottle in the right hand of m, a person already
// animated, if the person drinks.
func (b *Baker) holdBottle(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Drink) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	return human.Hold(m, building.HandBottle.Model(palette), human.BottleGrip)
}
