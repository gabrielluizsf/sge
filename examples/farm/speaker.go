package farm

import "github.com/gabrielluizsf/sge/vector"

// speaker plays the farm's sounds by name: at, when set, is where in the
// world it is heard from; nil is everywhere, as the HUD's are.
type speaker interface {
	play(name string, at *vector.Vector3)
}

// sound plays name from at, or everywhere for nil.
func (f *farm) sound(name string, at *vector.Vector3) {
	if f.speaker != nil {
		f.speaker.play(name, at)
	}
}
