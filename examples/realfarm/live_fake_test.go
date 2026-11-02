package main

import (
	"sync"
	"time"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// fakeMaker makes people that are a flat square, slowly if it is told to,
// and counts what it was asked: a stand-in for the real, slow, maker.
type fakeMaker struct {
	mu    sync.Mutex
	made  []string // the shapeKey of each person made
	plans int
	delay time.Duration
	gate  chan struct{} // if not nil, a make waits for it to be closed
}

func (f *fakeMaker) Make(c kit.Character) (person, error) {
	if f.gate != nil {
		<-f.gate
	}
	time.Sleep(f.delay)
	f.mu.Lock()
	defer f.mu.Unlock()
	f.made = append(f.made, shapeKey(c))
	return &fakePerson{f: f, model: square(c.Look.SkinTone)}, nil
}

func (f *fakeMaker) makes() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.made)
}

// fakePerson is a person whose colour is its skin tone, a grey.
type fakePerson struct {
	f     *fakeMaker
	model *gltf.Model
}

func (p *fakePerson) Model() *gltf.Model { return p.model }

func (p *fakePerson) Plan(c kit.Character) func() []int {
	p.f.mu.Lock()
	p.f.plans++
	p.f.mu.Unlock()
	return func() []int {
		g := c.Look.SkinTone
		p.model.Materials[0].BaseColor = [4]float64{g, g, g, 1}
		return []int{0}
	}
}

// square is a 1.7 metre tall square of grey g, facing the viewer.
func square(g float64) *gltf.Model {
	p := gltf.Primitive{
		Positions: [][3]float32{{-0.85, 0, 0}, {0.85, 0, 0}, {0.85, 1.7, 0}, {-0.85, 1.7, 0}},
		Normals:   [][3]float32{{0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}},
		UVs:       make([][2]float32, 4), Indices: []uint32{0, 1, 2, 0, 2, 3},
	}
	return &gltf.Model{
		Roots:     []int{0},
		Nodes:     []gltf.Node{{Mesh: 0, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}},
		Meshes:    []gltf.Mesh{{Primitives: []gltf.Primitive{p}}},
		Materials: []gltf.Material{{BaseColor: [4]float64{g, g, g, 1}, Texture: -1, MetalRough: -1, Normal: -1, Occlusion: -1, Roughness: 1, AlphaMode: "OPAQUE"}},
	}
}

// someone is a character by seed.
func someone(seed uint32) kit.Character {
	b, l := human.Random(human.Spec{Seed: seed})
	return kit.Character{Body: b, Look: l}
}

// withTone is c with the skin tone g, a fake person's grey.
func withTone(c kit.Character, g float64) kit.Character {
	c.Look.SkinTone = g
	return c
}
