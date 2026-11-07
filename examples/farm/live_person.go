package main

import (
	"github.com/gabrielluizsf/sge/examples/farm/kit"
	"github.com/gabrielluizsf/sge/gltf"
)

// person is a person made for the screen, whose colours change without
// making the mesh again.
type person interface {
	// Model is the person, standing.
	Model() *gltf.Model
	// Plan works out, off the frame, the change to c's colours; the
	// function it returns puts it on the model, quickly, on the frame, and
	// says which of the model's materials it changed. Plans are applied in
	// the order they were made.
	Plan(c kit.Character) (apply func() []int)
}

// maker makes people from nothing: slow.
type maker interface {
	Make(c kit.Character) (person, error)
}
