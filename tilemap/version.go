package tilemap

// Version changes whenever a cell comes to show another tile, which is the
// only way the map draws anything but the same (see game.Versioned): the
// game looks at where the map is, and at the camera, itself.
func (d drawing) Version() uint64 { return d.m.version }
