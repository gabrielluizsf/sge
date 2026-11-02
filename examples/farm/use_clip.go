package farm

// useClip is the clip a use of what the slot holds is shown by, as the
// suffix of its name, and how long it takes: a stroke of the axe, of the
// hammer where the axe is used on the anvil, or the act of anything else.
func useClip(slot int, anvil bool) (string, float64) {
	switch {
	case slot == slotAxe && anvil:
		return "_forge", ForgeSeconds
	case slot == slotAxe:
		return "_chop", actFor
	}
	return "_act", actFor
}
