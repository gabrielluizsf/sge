package terrain

// Area is a named kind of place laid over the world, such as water or
// woods: the cells In it, and the Tile each of them shows on the area's
// Layer. Layers are drawn from 0 up, the ground first, so an area of
// tiles transparent in places, as a lake's shores are, shows the layers
// under it there.
type Area struct {
	Name  string
	Layer int
	In    func(col, row int) bool
	Tile  func(col, row int) int
}
