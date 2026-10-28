package scene

// Forget lets go of scene name as it was left, so it is built anew the
// next time it is played: a room whose people come and go.
func (s *Stage) Forget(name string) { delete(s.saved, name) }
