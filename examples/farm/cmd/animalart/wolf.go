package main

import (
	"os"

	"github.com/gabrielluizsf/sge/animal"
)

// theWolf is the wolf's bones, measured on igor-lir's model: it faces -X
// as it is made.
var theWolf = beast{
	forward: animal.Vec{-1, 0, 0},
	body: []animal.Bone{
		{Name: "root", Head: animal.Vec{0, 3.8, -3.2}, Tail: animal.Vec{0, 3.8, -1.0}},
		{Name: "chest", Parent: "root", Head: animal.Vec{0, 3.8, -1.0}, Tail: animal.Vec{0, 4.0, 1.6}},
		{Name: "neck", Parent: "chest", Head: animal.Vec{0, 4.3, 1.6}, Tail: animal.Vec{0, 5.2, 2.7}},
		{Name: "head", Parent: "neck", Head: animal.Vec{0, 5.2, 2.7}, Tail: animal.Vec{0, 4.6, 4.2}},
		{Name: "tail", Parent: "root", Head: animal.Vec{0, 3.4, -3.8}, Tail: animal.Vec{0, 3.0, -4.9}},
		{Name: "rear", Parent: "root", Head: animal.Vec{0, 3.6, -3.9}, Tail: animal.Vec{0, 3.6, -3.2}},
	},
	feet:   [4]animal.Vec{{-0.45, 0, 1.48}, {0.44, 0, 1.50}, {-0.47, 0, -3.0}, {0.46, 0, -2.9}},
	hipY:   2.75,
	knee:   [2]float64{0.7, 0.7},
	chest:  "chest",
	gait:   animal.Gait{Stride: 2.4, Period: 0.55, Duty: 0.5, Lift: 0.6},
	crouch: 0,
}

func buildWolf(root *os.Root) ([]coat, error) {
	m, err := open(root, "animated_wolf")
	if err != nil {
		return nil, err
	}
	r, legs, err := theWolf.rig(m)
	if err != nil {
		return nil, err
	}
	walk := theWolf.walk(r, legs, 0.1333)
	return []coat{{name: "grey", model: m, motions: []motion{
		{name: "stand", times: []float64{0}},
		walk,
		wolfBite(r, legs),
		wolfFall(r, legs),
	}}}, nil
}
