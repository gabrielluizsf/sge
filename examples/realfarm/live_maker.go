package main

import (
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// bakerMaker makes people with the baker's human kit.
type bakerMaker struct{ b *bakes.Baker }

func (m bakerMaker) Make(c kit.Character) (person, error) {
	p, err := m.b.Person(c)
	if err != nil {
		return nil, err
	}
	return humanPerson{p}, nil
}

// humanPerson is a human.Person for the screen.
type humanPerson struct{ p *human.Person }

func (h humanPerson) Model() *gltf.Model { return h.p.Model }

func (h humanPerson) Plan(c kit.Character) func() []int {
	r := h.p.Plan(c.Look)
	return func() []int { return h.p.Apply(r) }
}
