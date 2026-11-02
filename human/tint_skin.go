package human

// skinStops are skin's average colour, sRGB, at tones 0, a quarter, a
// half, three quarters and 1: from the palest to the deepest.
var skinStops = [5]Colour{
	RGB(0.93, 0.77, 0.67), RGB(0.85, 0.65, 0.51), RGB(0.70, 0.50, 0.36), RGB(0.49, 0.32, 0.22), RGB(0.27, 0.17, 0.12),
}

// skinColour is the average colour of skin of tone, 0 the lightest to 1
// the darkest: the stops, each blended into the next.
func skinColour(tone float64) Colour {
	f := min(max(tone, 0), 1) * 4
	i := min(int(f), 3)
	a, b, t := skinStops[i], skinStops[i+1], f-float64(i)
	return RGB(a.R+(b.R-a.R)*t, a.G+(b.G-a.G)*t, a.B+(b.B-a.B)*t)
}

// skinKeep is how much of the skin's own hues stay when its tone is
// asked: its flush and veins, at a third.
const skinKeep = 0.33

// skinTint is the tint of skin of tone.
func skinTint(tone float64) tint {
	return tint{to: skinColour(tone).linear(), keep: skinKeep}
}
