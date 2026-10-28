// Package animation plays sprites one after another: an Animation is the
// frames and their pace, an Animator holds several by name and plays one.
package animation

import "github.com/gabrielluizsf/sge/sprite"

// Animation is frames shown FPS a second, over and over when Loop is set
// and stopping on the last otherwise.
type Animation struct {
	Frames []*sprite.Sprite
	FPS    float64
	Loop   bool
}
