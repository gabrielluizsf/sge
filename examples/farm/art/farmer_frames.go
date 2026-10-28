package art

// FarmerFrames is every frame of the farmer.
func FarmerFrames() map[string]string { return PersonFrames(FarmerLook) }

// PersonFrames is every frame of a person of look l, by direction and
// what it shows: standing, the two strides of a walk, pouring the can,
// raising the axe and bringing it down, and the two strokes of swimming. Left is right turned round.
func PersonFrames(l Look) map[string]string {
	f := map[string]string{
		"down":       front(l, pose{}, false),
		"down1":      front(l, pose{lift: [2]float64{1.4, 0}, swing: 0.7}, false),
		"down2":      front(l, pose{lift: [2]float64{0, 1.4}, swing: -0.7}, false),
		"downAct":    front(l, pose{can: true}, false),
		"up":         front(l, pose{}, true),
		"up1":        front(l, pose{lift: [2]float64{1.4, 0}, swing: 0.7}, true),
		"up2":        front(l, pose{lift: [2]float64{0, 1.4}, swing: -0.7}, true),
		"upAct":      front(l, pose{can: true}, true),
		"right":      side(l, pose{}),
		"right1":     side(l, pose{stride: 1, swing: 0.8}),
		"right2":     side(l, pose{stride: -1, swing: -0.8}),
		"rightAct":   side(l, pose{can: true}),
		"downChop":   front(l, pose{axe: windUp}, false),
		"downChop1":  front(l, pose{axe: strike}, false),
		"upChop":     front(l, pose{axe: windUp}, true),
		"upChop1":    front(l, pose{axe: strike}, true),
		"rightChop":  side(l, pose{axe: windUp}),
		"rightChop1": side(l, pose{axe: strike}),
	}
	for _, s := range []string{"", "1", "2", "Act", "Chop", "Chop1"} {
		f["left"+s] = mirrored(f["right"+s], FarmerW)
	}
	swimFrames(f)
	return f
}

// mirrored is the drawing svg, w units wide, turned round left to right.
func mirrored(svg string, w float64) string {
	i := indexAfter(svg, "</defs>\n")
	return svg[:i] + `<g transform="` + "matrix(-1 0 0 1 " + ftoa(w) + ` 0)">` + "\n" + svg[i:len(svg)-len("</svg>\n")] + "</g>\n</svg>\n"
}
