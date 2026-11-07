package art

import "fmt"

// lying is the farmer's drawing svg laid on their back, head to the left,
// across the lower part of the picture where a bed's top is, the chest
// swelling by swell: a rise of the breath.
func lying(svg string, swell float64) string {
	const length = 0.57 // the figure shrunk to fit the picture's width
	i := indexAfter(svg, "</defs>\n")
	return svg[:i] + fmt.Sprintf(`<g transform="translate(0 34) rotate(-90) scale(%.3f %.3f)">`+"\n", length*swell, length) +
		svg[i:len(svg)-len("</svg>\n")] + "</g>\n</svg>\n"
}

// sleepFrames are the farmer asleep, whichever way they face: two frames
// of a slow breath.
func sleepFrames(f map[string]string) {
	f["Sleep"] = lying(f["down"], 1)
	f["Sleep1"] = lying(f["down"], 1.06)
}
