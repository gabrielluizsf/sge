package main

import (
	"math"
	"slices"

	"github.com/gabrielluizsf/sge/human"
)

// cycle is the item d (+1 or -1) places from cur in list, round the ends;
// the first or last when cur is not in it. "" for an empty list.
func cycle(list []string, cur string, d int) string {
	if len(list) == 0 {
		return ""
	}
	i := slices.Index(list, cur)
	if i < 0 {
		if d > 0 {
			return list[0]
		}
		return list[len(list)-1]
	}
	return list[(i+d+len(list))%len(list)]
}

// cycleHat is cycle over list with "" (no hat) before its first item.
func cycleHat(list []string, cur string, d int) string {
	return cycle(append([]string{""}, list...), cur, d)
}

// nearest is the index of the colour of list nearest c, -1 when c is
// not set or list is empty.
func nearest(list []human.Colour, c human.Colour) int {
	best, at := math.Inf(1), -1
	if !c.Set {
		return -1
	}
	for i, x := range list {
		if d := math.Abs(x.R-c.R) + math.Abs(x.G-c.G) + math.Abs(x.B-c.B); d < best {
			best, at = d, i
		}
	}
	return at
}

// cycleColour is the colour d places from c in list, round the ends. The
// unset colour (the pack's own) is a stop of its own before the first, if
// original is true.
func cycleColour(list []human.Colour, c human.Colour, d int, original bool) human.Colour {
	n := len(list)
	if original {
		n++
	}
	if n == 0 {
		return c
	}
	i := nearest(list, c)
	if original {
		i++ // the original is stop 0
	} else if i < 0 {
		i = 0
	}
	i = (i + d + n) % n
	if original {
		if i == 0 {
			return human.Colour{}
		}
		i--
	}
	return list[i]
}
