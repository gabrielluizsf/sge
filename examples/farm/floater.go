package main

import "github.com/gabrielluizsf/sge/vector"

// floater is a word rising from where something happened, for a second.
type floater struct {
	at   vector.Vector3
	text string
	ink  rune
	age  float64
}

const floatFor = 1.2

// float shows text rising from at, in the palette colour ink.
func (f *farm) float(at vector.Vector3, text string, ink rune) {
	f.floats = append(f.floats, &floater{at: at, text: text, ink: ink})
}

// drift moves the floaters up, and lets the old ones go.
func (f *farm) drift(dt float64) {
	kept := f.floats[:0]
	for _, w := range f.floats {
		if w.age += dt; w.age < floatFor {
			w.at.Y -= 24 * dt
			kept = append(kept, w)
		}
	}
	f.floats = kept
}
