package art

// FarmerFrames is every frame of the farmer, by direction and what it
// shows: standing, the two strides of a walk, pouring the can, and the
// two strokes of swimming. Left is right turned round.
func FarmerFrames() map[string]string {
	f := map[string]string{
		"down":     front(pose{}, false),
		"down1":    front(pose{lift: [2]float64{1.4, 0}, swing: 0.7}, false),
		"down2":    front(pose{lift: [2]float64{0, 1.4}, swing: -0.7}, false),
		"downAct":  front(pose{can: true}, false),
		"up":       front(pose{}, true),
		"up1":      front(pose{lift: [2]float64{1.4, 0}, swing: 0.7}, true),
		"up2":      front(pose{lift: [2]float64{0, 1.4}, swing: -0.7}, true),
		"upAct":    front(pose{can: true}, true),
		"right":    side(0, 0, false),
		"right1":   side(1, 0.8, false),
		"right2":   side(-1, -0.8, false),
		"rightAct": side(0, 0, true),
	}
	for _, s := range []string{"", "1", "2", "Act"} {
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
