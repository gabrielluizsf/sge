package art

import (
	"strconv"
	"strings"
)

// indexAfter is the index just past sep's first place in s.
func indexAfter(s, sep string) int { return strings.Index(s, sep) + len(sep) }

func ftoa(f float64) string { return strconv.FormatFloat(f, 'g', -1, 64) }
