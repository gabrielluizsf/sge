package kit

import (
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/save"
	"github.com/gabrielluizsf/sge/textfield"
)

// The character's file in a game's save folder, its format's version, and
// the most bytes it may have.
const (
	characterFile    = "look"
	characterVersion = 1
	characterLimit   = 16 << 10
)

// CharacterStore is where the game titled title keeps its character.
func CharacterStore(title string) (*save.Store, error) { return save.Open(title, characterVersion) }

// SaveCharacter writes c into s.
func SaveCharacter(s *save.Store, c Character) error { return s.Write(characterFile, c) }

// LoadCharacter reads the character saved in s: ok is false when there is
// none, or the file is not a good one (unreadable, too large, unknown
// fields, a file or number the creator would not make): the default
// player is then kept.
func LoadCharacter(s *save.Store) (c Character, ok bool) {
	if !s.Has(characterFile) {
		return c, false
	}
	if err := s.Read(characterFile, &c, characterLimit); err != nil {
		return Character{}, false
	}
	if !textfield.Valid(c.Name, MaxName) {
		c.Name = "" // a bad name is the default one, not a reason to lose the look
	}
	c.Look = human.Mended(c.Look)
	if c.Check() != nil {
		return Character{}, false
	}
	return c, true
}
