package game

import "github.com/gabrielluizsf/antui/canvas"

// Script is behaviour attached to an object, run every frame with the
// seconds since the last.
type Script interface {
	Update(o *Object, dt float64)
}

// Starter is a script that wants to know when its object first joins a
// running world, before its first Update.
type Starter interface {
	Start(o *Object)
}

// CollisionHandler is a script told when its object touches another.
type CollisionHandler interface {
	Collision(o, other *Object)
}

// ScriptFunc is a function that is a Script.
type ScriptFunc func(o *Object, dt float64)

// Update calls f.
func (f ScriptFunc) Update(o *Object, dt float64) { f(o, dt) }

// Painter is a script that draws its object itself, after the object's
// sprite and in the same order: text, a health bar, anything a sprite is
// not.
type Painter interface {
	Paint(o *Object, cv *canvas.Canvas)
}
