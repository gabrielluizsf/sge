package art

import "fmt"

// The farmer swimming shows from the waist up, sunk to it: the water cuts
// the drawing at waist and moves it sink lower, so its line is just above
// the picture's bottom, where the farmer stands.
const waist, sink = 27, 10.5

// swimming is the farmer's drawing svg sunk to the waist in a ring of
// ripples, bobbing bob units lower, the ring spread wider by ring.
func swimming(svg string, bob, ring float64) string {
	i := indexAfter(svg, "<defs>\n")
	j := indexAfter(svg, "</defs>\n")
	y := waist + sink + bob
	rx := 7.5 + ring
	ripple := func(sweep int) string {
		return fmt.Sprintf(`<path d="M%.2f %.2f a%.2f 2.4 0 0 %d %.2f 0" fill="none" stroke="%s" stroke-width="0.9" stroke-linecap="round" opacity="0.9"/>`+"\n",
			12-rx, y, rx, sweep, 2*rx, foam)
	}
	return svg[:i] + fmt.Sprintf(`<clipPath id="waist"><rect x="0" y="0" width="%d" height="%d"/></clipPath>`+"\n", FarmerW, waist) +
		svg[i:j] +
		fmt.Sprintf(`<ellipse cx="12" cy="%.2f" rx="%.2f" ry="2.4" fill="%s" opacity="0.55"/>`+"\n", y, rx, waterLight) +
		ripple(1) +
		fmt.Sprintf(`<g transform="translate(0 %.2f)"><g clip-path="url(#waist)">`+"\n", sink+bob) +
		svg[j:len(svg)-len("</svg>\n")] + "</g></g>\n" +
		ripple(0) + "</svg>\n"
}

// swimFrames are the farmer's frames in the water, by direction: two
// strokes, the arms of one stride and then the other, bobbing between.
func swimFrames(f map[string]string) {
	for _, dir := range []string{"down", "up", "right"} {
		f[dir+"Swim"] = swimming(f[dir+"1"], 0, 0)
		f[dir+"Swim1"] = swimming(f[dir+"2"], 0.7, 1.2)
	}
	for _, s := range []string{"Swim", "Swim1"} {
		f["left"+s] = mirrored(f["right"+s], FarmerW)
	}
}
