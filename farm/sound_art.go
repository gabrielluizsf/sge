package farm

import "github.com/gabrielluizsf/sge/sound"

// SoundArt is what an Art may also have that brings sounds of its own: by
// name, the versions of each play instead of those package sfx makes up,
// and every other sound stays as it is.
type SoundArt interface {
	Sounds() map[string][]*sound.Clip
}

// useSounds has s play the sounds of a in place of its own, if it has any.
func (s *mixerSpeaker) useSounds(a Art) {
	sa, ok := a.(SoundArt)
	if !ok {
		return
	}
	for name, versions := range sa.Sounds() {
		if len(versions) > 0 {
			s.clips[name] = versions
		}
	}
}
