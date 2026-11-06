package farm

import "github.com/gabrielluizsf/antui/canvas"

// pictureKeep is how many pictures each of the two generations of the
// picture cache holds.
const pictureKeep = 1024

// pictureCache holds the drawings rendered lately, by source and size, in
// two generations: a picture asked for is found in either and moves to
// the newer; when that is full the older is dropped. What every sprite on
// screen holds stays with it, so a picture dropped from here is only
// drawn again, the same, if a new sprite asks for it: a walk past a
// thousand villagers does not keep every frame of them for good.
type pictureCache struct {
	newer, older map[string]*canvas.Canvas
}

func (c *pictureCache) get(key string) (*canvas.Canvas, bool) {
	if cv, ok := c.newer[key]; ok {
		return cv, true
	}
	cv, ok := c.older[key]
	if ok {
		c.put(key, cv)
	}
	return cv, ok
}

func (c *pictureCache) put(key string, cv *canvas.Canvas) {
	if c.newer == nil || len(c.newer) >= pictureKeep {
		c.older, c.newer = c.newer, map[string]*canvas.Canvas{}
	}
	c.newer[key] = cv
}

// size is how many pictures the cache holds.
func (c *pictureCache) size() int { return len(c.newer) + len(c.older) }
