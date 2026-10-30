package farm

// Keyboard is the player at the keys and the mouse, as a window gives
// them.
func Keyboard() Controller { return keyboard{} }

// Bot is the demo's player, who plays the farm by itself and keeps the
// first tends plots of the field.
func Bot(tends int) Controller { return &bot{tends: tends} }
