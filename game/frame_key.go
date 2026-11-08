package game

import (
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/transform"
	"github.com/gabrielluizsf/sge/vector"
)

// frameKey is all that Draw puts on a canvas depends on, bar what the
// painters of the world keep to themselves (see Versioned): when the keys of
// two frames are the same, so are their pixels.
type frameKey struct {
	valid         bool
	canvas        *canvas.Canvas
	width, height int
	camera        vector.Vector3
	background    canvas.Color
	covered       bool
	projection    [4]vector.Vector3 // where the origin and the unit steps go
	objects       []objectLook
	versions      []uint64 // of the painters, in the order of the objects
}

// objectLook is what the game sees of an object that is drawn.
type objectLook struct {
	object    *Object
	transform transform.Transform
	footing   vector.Vector3
	z         int
	screen    bool
	look      *sprite.Sprite
	image     *canvas.Canvas
	pivot     [2]float64
}

// fill sets k to the key of drawing the world onto cv now, and reports
// whether every painter of it has a version.
func (k *frameKey) fill(g *Game, cv *canvas.Canvas) (steady bool) {
	k.valid, k.canvas, k.width, k.height = true, cv, cv.Width, cv.Height
	k.camera, k.background, k.covered = g.camera(), g.Background, g.Covered
	k.projection = [4]vector.Vector3{
		g.project(vector.Vector3{}), g.project(vector.Vector3{X: 1}),
		g.project(vector.Vector3{Y: 1}), g.project(vector.Vector3{Z: 1}),
	}
	k.objects, k.versions = k.objects[:0], k.versions[:0]
	steady = true
	for _, o := range g.Objects() {
		if o.Hidden || o.Look() == nil && !o.draws() {
			continue
		}
		k.objects = append(k.objects, lookOf(o))
		var ok bool
		k.versions, ok = o.versions(k.versions)
		steady = steady && ok
	}
	return steady
}

// lookOf is what the game sees of o.
func lookOf(o *Object) objectLook {
	look := o.Look()
	l := objectLook{object: o, transform: o.Transform, footing: o.Footing, z: o.Z, screen: o.Screen, look: look}
	if look != nil {
		l.image, l.pivot = look.Image, [2]float64{look.PivotX, look.PivotY}
	}
	return l
}

// same reports whether k and o are the same key.
func (k *frameKey) same(o *frameKey) bool {
	return k.valid && o.valid && k.canvas == o.canvas && k.width == o.width && k.height == o.height &&
		k.camera == o.camera && k.background == o.background && k.covered == o.covered &&
		k.projection == o.projection && slices.Equal(k.objects, o.objects) && slices.Equal(k.versions, o.versions)
}
