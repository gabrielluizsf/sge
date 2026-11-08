// Command glview opens a window and draws in it, through render3d/gl,
// a textured cube spinning and a skinned arm bending, on a ground in the
// sun's shadow: for a person to look at what the OpenGL backend draws.
// Escape, or closing the window, quits; -seconds quits by itself.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/gl"
)

func main() {
	seconds := flag.Float64("seconds", 0, "quit after this many seconds (0: never)")
	flag.Parse()
	if err := run(*seconds); err != nil {
		fmt.Fprintln(os.Stderr, "glview:", err)
		os.Exit(1)
	}
}

// run shows the view until the window closes, or for seconds when not 0,
// and says how many frames a second it drew.
func run(seconds float64) error {
	win, err := antui.Open("sge render3d/gl", 960, 540)
	if err != nil {
		return err
	}
	defer win.Close()
	dev, err := gl.Window(win)
	if err != nil {
		return err
	}
	defer dev.Close()
	r := render3d.NewRenderer(dev)
	v, err := newView(r)
	if err != nil {
		return err
	}
	start, frames := time.Now(), 0
	for win.Begin() && !win.KeyDown(antui.KeyEscape) {
		t := time.Since(start).Seconds()
		if seconds > 0 && t > seconds {
			break
		}
		if w, h := dev.Size(); w != win.Width() || h != win.Height() {
			dev.Resize(win.Width(), win.Height())
		}
		v.at(t, [2]int{win.Width(), win.Height()})
		if err := r.Draw(&v.scene); err != nil {
			return err
		}
		dev.Present()
		win.End()
		frames++
	}
	fmt.Printf("%d frames, %.1f a second\n", frames, float64(frames)/time.Since(start).Seconds())
	return nil
}
