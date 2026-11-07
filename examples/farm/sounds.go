package main

import (
	"embed"
	"encoding/binary"
	"strings"

	"github.com/gabrielluizsf/sge/sound"
)

// soundFiles are the game's own sounds, as sounds/SOURCE.txt says, each named
// by the sound it is and which of its versions.
//
//go:embed sounds/*.pcm
var soundFiles embed.FS

// Sounds are the look's sounds by name, in the order of their versions.
// Those it has none of, the water's and the birds', stay as package sfx makes
// them.
func (a *realArt) Sounds() map[string][]*sound.Clip {
	files, err := soundFiles.ReadDir("sounds")
	if err != nil {
		return nil
	}
	out := map[string][]*sound.Clip{}
	for _, f := range files {
		raw, err := soundFiles.ReadFile("sounds/" + f.Name())
		if err != nil {
			continue
		}
		name := soundName(f.Name())
		out[name] = append(out[name], clipOf(raw))
	}
	return out
}

// soundName is the sound file holds a version of: "step_grass_2.pcm" is of
// "step_grass".
func soundName(file string) string {
	return file[:strings.LastIndex(file, "_")]
}

// clipOf is the clip of raw, mono 16-bit little-endian samples.
func clipOf(raw []byte) *sound.Clip {
	c := &sound.Clip{Rate: sound.Rate, Samples: make([]float32, len(raw)/2)}
	for i := range c.Samples {
		c.Samples[i] = float32(int16(binary.LittleEndian.Uint16(raw[2*i:]))) / 32768
	}
	return c
}
