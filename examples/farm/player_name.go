package farm

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/textfield"
)

// MaxName is the most characters of the player's name.
const MaxName = 20

// NamedArt is an Art that knows what the player called their farmer, from
// its character creation screen: "" for none, and the farmer then has the
// founder's name. A name that is not a good one (see textfield.Valid) is
// not taken.
type NamedArt interface {
	PlayerName() string
}

// playerName is the name the art gives the player, or the founder's.
func playerName(a Art) string {
	if n, ok := a.(NamedArt); ok {
		if s := n.PlayerName(); textfield.Valid(s, MaxName) {
			return s
		}
	}
	return founderName
}

// RandomName is a name from the villagers' lists, a woman's or a man's,
// by seed: the same seed, the same name.
func RandomName(woman bool, seed uint32) string {
	names := menNames
	if woman {
		names = womenNames
	}
	return names[random.New(seed).Intn(len(names))]
}
