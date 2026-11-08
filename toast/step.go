package toast

import "math"

// Step moves time on by dt: toasts age, ease toward their rows and, once
// faded out, are dropped. It is the list to keep.
func Step(list []*Toast, dt float64) []*Toast {
	kept := list[:0]
	for _, t := range list {
		t.Age += dt
		if !t.Gone() {
			kept = append(kept, t)
		}
	}
	place(kept)
	k := 1 - math.Exp(-Smooth*dt)
	for _, t := range kept {
		t.Row += (t.Target - t.Row) * k
	}
	return kept
}
