package farm

import "github.com/gabrielluizsf/sge/vector"

// cue is a sound due to play in a moment, from a place: a blow's, as it
// lands, not as it is raised.
type cue struct {
	name string
	at   vector.Vector3
	in   float64
}

// soundIn plays c's sound when its moment comes.
func (f *farm) soundIn(c cue) { f.cues = append(f.cues, c) }

// ring lets dt seconds go by on the cues, playing those now due.
func (f *farm) ring(dt float64) {
	kept := f.cues[:0]
	for _, c := range f.cues {
		if c.in -= dt; c.in > 0 {
			kept = append(kept, c)
			continue
		}
		f.sound(c.name, &c.at)
	}
	f.cues = kept
}
