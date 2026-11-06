package farm

import "github.com/gabrielluizsf/sge/sprite"

// bodyFirst is the frame a person is first seen in, standing facing us.
func bodyFirst(b Body) *sprite.Sprite { return b.Clip(AnimStand, "down").Frames[0] }

// bodyFrames are the frames a person is mostly seen in, standing and
// walking each way, standing first: what is drawn ahead of the person
// coming into view, after the first frame of everyone in the town.
func bodyFrames(b Body) []*sprite.Sprite {
	var frames []*sprite.Sprite
	for _, anim := range []string{AnimStand, AnimWalk} {
		for _, dir := range Facings(b.Facings()) {
			frames = append(frames, b.Clip(anim, dir).Frames...)
		}
	}
	return frames
}
