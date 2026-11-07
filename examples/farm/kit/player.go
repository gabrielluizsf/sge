package kit

// playerTag names the character the player made, "" while the player is
// the one the game was made with.
var playerTag string

// playerName is the name the player gave their character, "" for none.
var playerName string

// PlayerName is the name the player gave their character, "" if they
// gave none: the game then keeps its own.
func PlayerName() string { return playerName }

// PlayerTag is what a task name of villager v adds to tell the character
// the player made: " look <tag>" if v is the player and made one, else "".
func PlayerTag(v Villager) string {
	if v.ID == "player" && playerTag != "" {
		return " look " + playerTag
	}
	return ""
}

// SpriteTag is what the name of a sprite of villager id adds to tell the
// character the player made, so each look keeps its own pictures:
// "-<tag>" for the player who made one, else "".
func SpriteTag(id string) string {
	if id == "player" && playerTag != "" {
		return "-" + playerTag
	}
	return ""
}

// SetPlayer makes the player villager the character c is: their body and
// look, and a tag that names the pictures baked of them. It must not be
// called while Villagers are read elsewhere; c must be Checked.
func SetPlayer(c Character) {
	for i := range Villagers {
		if Villagers[i].ID == "player" {
			Villagers[i].Body, Villagers[i].Look = c.Body, c.Look
		}
	}
	playerTag, playerName = c.Tag(), c.Name
}

// PlayerCharacter is the player villager as a character.
func PlayerCharacter() Character {
	for _, v := range Villagers {
		if v.ID == "player" {
			return Character{Body: v.Body, Look: v.Look, Name: playerName}
		}
	}
	return Character{}
}

// defaultPlayer is the player the game was made with.
var defaultPlayer = PlayerCharacter()

// ResetPlayer makes the player the one the game was made with.
func ResetPlayer() {
	SetPlayer(defaultPlayer)
	playerTag, playerName = "", ""
}
