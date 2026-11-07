package farm

// sleepAndSave puts the player to bed and, the night over, saves the game
// where the pause menu's button would: what is lost to a crash or a closed
// window is at most a day. A game with no saving, a demo, only sleeps.
func (f *farm) sleepAndSave() bool {
	if !f.sleep() {
		return false
	}
	if f.ui.save != nil {
		f.ui.save()
	}
	return true
}
