package human

// The faces a game asks for most, each as far as w, from 0, the face at
// rest, to 1, as MakeHuman's face pose units make them. Put more than one
// together with Blend.

// Smile is both corners of the mouth drawn up, and the cheeks with them.
func Smile(w float64) Expression {
	return Expression{"MouthLeftPullUp": w, "MouthRightPullUp": w, "LeftCheekUp": 0.5 * w, "RightCheekUp": 0.5 * w}
}

// Sad is the corners of the mouth down, and the inner brows raised.
func Sad(w float64) Expression {
	return Expression{"MouthLeftPullDown": w, "MouthRightPullDown": w, "LeftInnerBrowUp": 0.8 * w, "RightInnerBrowUp": 0.8 * w}
}

// Angry is the brows down, the nose wrinkled and the lip drawn up.
func Angry(w float64) Expression {
	return Expression{"LeftBrowDown": w, "RightBrowDown": w, "NoseWrinkler": 0.6 * w, "UpperLipUp": 0.4 * w}
}

// Surprised is the brows and lids up, and the jaw dropped a little.
func Surprised(w float64) Expression {
	return Expression{"LeftOuterBrowUp": w, "RightOuterBrowUp": w, "LeftInnerBrowUp": w, "RightInnerBrowUp": w,
		"LeftUpperLidOpen": w, "RightUpperLidOpen": w, "JawDrop": 0.35 * w}
}

// Pain is the eyes shut tight, the brows down and the teeth bared.
func Pain(w float64) Expression {
	return Expression{"LeftUpperLidClosed": w, "RightUpperLidClosed": w, "LeftBrowDown": 0.8 * w, "RightBrowDown": 0.8 * w,
		"MouthLeftPullSide": 0.7 * w, "MouthRightPullSide": 0.7 * w}
}

// Wink is the left eye shut, with a little of a smile on that side.
func Wink(w float64) Expression {
	return Expression{"LeftUpperLidClosed": w, "LeftCheekUp": 0.6 * w, "MouthLeftPullUp": 0.4 * w}
}

// Blend is the faces together: of a unit that more than one moves, the most.
func Blend(faces ...Expression) Expression {
	out := Expression{}
	for _, face := range faces {
		for unit, w := range face {
			out[unit] = max(out[unit], w)
		}
	}
	return out
}
