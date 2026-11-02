package main

import (
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/gl"
	"github.com/gabrielluizsf/sge/render3d/soft"
)

// deviceMaker makes the device a stage draws on, w by h, and how to let
// it go.
type deviceMaker func(w, h int) (render3d.Device, func(), error)

// bestDevice is OpenGL, offscreen, where the machine has it, and the
// software device, which every machine has, where it has not.
func bestDevice(w, h int) (render3d.Device, func(), error) {
	if d, err := gl.Offscreen(w, h); err == nil {
		return d, func() { d.Close() }, nil
	}
	return softDevice(w, h)
}

// softDevice is the software device.
func softDevice(w, h int) (render3d.Device, func(), error) {
	return soft.New(w, h), func() {}, nil
}
