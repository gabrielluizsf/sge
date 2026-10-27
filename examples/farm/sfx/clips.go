package sfx

import "github.com/gabrielluizsf/sge/sound"

// Clips are the farm's sounds played once, by name, each with versions
// to play in turn, so a step is never twice the same.
func Clips() map[string][]*sound.Clip {
	one := func(c *sound.Clip) []*sound.Clip { return []*sound.Clip{c} }
	return map[string][]*sound.Clip{
		"step_grass": versions(4, grassStep),
		"step_soil":  versions(4, soilStep),
		"step_path":  versions(4, pathStep),
		"stroke":     versions(3, stroke),
		"splash":     versions(1, splash),
		"chirp":      versions(5, chirp),
		"ui_tick":    one(knocks(1300)),
		"ui_open":    one(knocks(880, 1320)),
		"ui_close":   one(knocks(1320, 880)),
		"coin":       one(bells(0.06, 1568, 2093)),
		"sell":       one(bells(0.07, 1047, 1319, 1568, 2093)),
		"deny":       one(deny()),
		"plant":      one(plant()),
		"pour":       one(pour()),
		"harvest":    one(harvest()),
	}
}

// versions is n versions of a sound, each made from a seed of its own.
func versions(n int, make1 func(uint32) *sound.Clip) []*sound.Clip {
	c := make([]*sound.Clip, n)
	for i := range c {
		c[i] = make1(uint32(i*97 + 1))
	}
	return c
}

// Loops are the farm's ambience, by name: leaves near trees, water near
// the shore.
func Loops() map[string]*sound.Clip {
	return map[string]*sound.Clip{"rustle": rustle(), "lap": lap()}
}
