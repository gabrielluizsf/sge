package game

// Destroy takes the object out of the world at the end of this step.
func (o *Object) Destroy() { o.destroyed = true }
