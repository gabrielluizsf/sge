package scene

import "github.com/gabrielluizsf/sge/game"

// Switch leaves the scene playing and plays scene name: as it was left,
// or, the first time, as build sets it up on the game, adding its objects
// and placing its camera. The objects of the scene left leave the game at
// the end of the step, kept as they are, with the camera, background and
// Covered of the game; Keep's objects stay. Switching to the scene
// playing does nothing.
func (s *Stage) Switch(name string, build func(g *game.Game)) {
	if name == s.Current {
		return
	}
	g := s.Game
	if s.saved == nil {
		s.saved = map[string]*shot{}
	}
	left := &shot{camera: g.Camera, background: g.Background, covered: g.Covered}
	for _, o := range g.Objects() {
		if s.Keep != nil && s.Keep(o) {
			continue
		}
		left.objects = append(left.objects, o)
		o.Destroy()
	}
	s.saved[s.Current] = left
	s.Current = name
	back, ok := s.saved[name]
	if !ok {
		build(g)
		return
	}
	delete(s.saved, name)
	for _, o := range back.objects {
		g.Add(o)
	}
	g.Camera, g.Background, g.Covered = back.camera, back.background, back.covered
}

// Played reports whether scene name has been played and left, and so
// comes back as it was.
func (s *Stage) Played(name string) bool {
	_, ok := s.saved[name]
	return ok
}
