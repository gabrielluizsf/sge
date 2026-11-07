package main

import (
	"os"

	"github.com/gabrielluizsf/sge/animal"
)

// theSheep is the sheep's bones, measured on igor-lir's model: it faces
// -X as it is made, and its legs are posed in mid-stride.
var theSheep = beast{
	forward: animal.Vec{-1, 0, 0},
	body: []animal.Bone{
		{Name: "root", Head: animal.Vec{0, 2.8, -2.0}, Tail: animal.Vec{0, 2.8, -0.5}},
		{Name: "chest", Parent: "root", Head: animal.Vec{0, 2.8, -0.5}, Tail: animal.Vec{0, 2.9, 1.2}},
		{Name: "neck", Parent: "chest", Head: animal.Vec{0, 3.3, 1.2}, Tail: animal.Vec{0, 4.3, 2.4}},
		{Name: "head", Parent: "neck", Head: animal.Vec{0, 4.3, 2.4}, Tail: animal.Vec{0, 4.3, 3.5}},
		{Name: "tail", Parent: "root", Head: animal.Vec{0, 3.2, -3.0}, Tail: animal.Vec{0, 2.9, -3.6}},
		{Name: "rear", Parent: "root", Head: animal.Vec{0, 2.9, -3.2}, Tail: animal.Vec{0, 2.9, -2.0}},
	},
	feet:   [4]animal.Vec{{-0.52, 0, 1.23}, {0.70, 0, 0.48}, {-0.44, 0, -2.77}, {0.47, 0, -1.84}},
	hipY:   1.45,
	knee:   [2]float64{0.55, 0.55},
	chest:  "chest",
	gait:   animal.Gait{Stride: 1.6, Period: 0.9, Duty: 0.6, Lift: 0.35},
	crouch: 0,
}

func buildSheep(root *os.Root) ([]coat, error) {
	m, err := open(root, "animated_sheep")
	if err != nil {
		return nil, err
	}
	r, legs, err := theSheep.rig(m)
	if err != nil {
		return nil, err
	}
	walk := theSheep.walk(r, legs, 0.15)
	return []coat{{name: "white", model: m, motions: []motion{
		{name: "stand", times: []float64{0}},
		walk,
		graze(r, legs, "neck", "head", rad(38)),
	}}}, nil
}
