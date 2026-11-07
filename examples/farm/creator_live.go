package main

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// advance takes what the preview has ready, moves the camera dt seconds
// on, and draws the person if anything changed. It is cheap when nothing
// did, and never waits for a person being made.
func (c *creator) advance(dt float64) {
	if c.live == nil || c.stageAt.W <= 0 {
		return
	}
	in := c.inner()
	if u, ok := c.live.Take(); ok {
		if err := c.stage.Show(u, in.W, in.H); err != nil {
			c.shownErr = err
		}
	}
	if c.view.Step(dt) {
		c.stage.Dirty()
	}
}

// inner is where the picture goes, in the stage.
func (c *creator) inner() rect {
	r := c.stageAt
	return rect{r.X + stageInset, r.Y + stageInset, r.W - 2*stageInset, r.H - 2*stageInset}
}

// stageImage is the picture of the person to draw in r, nil if there is
// none yet.
func (c *creator) stageImage(r rect) *canvas.Canvas {
	if c.stage == nil {
		return nil
	}
	img, err := c.stage.Render(&c.view, r.W, r.H)
	if err != nil {
		c.shownErr = err
	}
	return img
}
