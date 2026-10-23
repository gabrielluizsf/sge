package game

// Game is the world the object is in, or nil before it is added.
func (o *Object) Game() *Game { return o.game }
