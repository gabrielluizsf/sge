package main

import (
	"os"

	"github.com/gabrielluizsf/sge/gltf"
)

// animalDef is an animal to bake: its name, how many of a cell a unit of
// its model is, and how to read the model's folder into its coats.
type animalDef struct {
	name  string
	unit  float64
	build func(root *os.Root) ([]coat, error)
}

// coat is a model of an animal in one colouring, rigged, and the motions
// to bake of it.
type coat struct {
	name    string
	model   *gltf.Model
	motions []motion
}

// motion is a clip of a coat, baked at the seconds of times, a frame at
// each; a clip of "" is the model as it stands. Speed is how fast the
// animal goes in the walk, in cells a second; Once is a motion that is
// run through, not looped.
type motion struct {
	name  string
	clip  string
	times []float64
	once  bool
	speed float64
}

// animals are those baked, in order.
var animals = []animalDef{
	{"cow", 2.6, buildCow},
	{"sheep", 0.15, buildSheep},
	{"wolf", 0.1333, buildWolf},
	{"horse", 0.7, buildHorse},
}

// open reads the model in the folder dir of root, with the reader's limits.
func open(root *os.Root, dir string) (*gltf.Model, error) {
	sub, err := root.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer sub.Close()
	return gltf.Open(sub, "scene.gltf", gltf.DefaultLimits)
}

// loop is n times evenly round a loop of period seconds.
func loop(n int, period float64) []float64 {
	t := make([]float64, n)
	for i := range t {
		t[i] = period * float64(i) / float64(n)
	}
	return t
}

// run is n times from the first to the last second of a motion run once.
func run(n int, length float64) []float64 {
	t := make([]float64, n)
	for i := range t {
		t[i] = length * float64(i) / float64(n-1) * (1 - 1e-9) // just short of the wrap to the start
	}
	return t
}
