package farm

import "github.com/gabrielluizsf/sge/toast"

// notify shows text as a toast at the top of the screen.
func (f *farm) notify(text string) { f.notices = toast.Add(f.notices, text) }

// fade moves the toasts on: they slide in, hold, and fade out.
func (f *farm) fade(dt float64) { f.notices = toast.Step(f.notices, dt) }
