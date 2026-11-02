package kit

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/gabrielluizsf/sge/human"
)

// Character is a person the player made: their body, how they look and
// their name.
type Character struct {
	Body human.Body
	Look human.Look
	// Name is what the player calls them, 1 to MaxName characters, or ""
	// for the game's own default; it is not part of the Tag.
	Name string `json:"Name,omitempty"`
}

// Tag is a short name for the character, to tell the pictures baked of
// them from any other's: the same character, the same tag.
func (c Character) Tag() string {
	b, _ := json.Marshal(Character{Body: c.Body, Look: c.Look})
	h := fnv.New64a()
	h.Write(b)
	return fmt.Sprintf("%012x", h.Sum64()&0xffffffffffff)
}
