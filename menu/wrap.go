package menu

import "strings"

// Wrap breaks text into lines no wider than w pixels, a word's width given
// by wide, keeping the text's own line breaks.
func Wrap(text string, w int, wide func(string) int) []string {
	var out []string
	space := wide("a a") - wide("aa")
	for _, para := range strings.Split(strings.ReplaceAll(text, "\r", ""), "\n") {
		line, used := "", 0
		for _, word := range strings.Fields(para) {
			ww := wide(word)
			if line != "" && used+space+ww > w {
				out = append(out, line)
				line, used = "", 0
			}
			if line != "" {
				line, used = line+" ", used+space
			}
			line, used = line+word, used+ww
		}
		out = append(out, line)
	}
	return out
}

// Widths is a wide function for Wrap that asks art for each word's width
// once, remembering it in the map it is given.
func Widths(art Art, known map[string]int) func(string) int {
	return func(word string) int {
		if n, ok := known[word]; ok {
			return n
		}
		known[word] = art.Text(word, 'k').Width
		return known[word]
	}
}
