package toast

import "strings"

// Life is how long a toast with text stays before it leaves: about two
// seconds for a short one, a little more for each word, at most three and a
// half.
func Life(text string) float64 {
	return min(3.5, max(2, 1.6+0.2*float64(len(strings.Fields(text)))))
}
