package terrain

// Edit sets cell col, row of layer to tile id for good, over what the
// areas give it, tilemap.Empty clearing it: a tree cut down, a field
// tilled. Areas stay as they are: In still reports the cell in the woods.
// A cell beyond the reach of keys is left alone.
func (g *Generator) Edit(layer, col, row, id int) {
	k, ok := editKey(layer, col, row)
	if !ok {
		return
	}
	if g.edits == nil {
		g.edits = map[int]int{}
	}
	g.edits[k] = id
}

// edited is the tile Edit set on cell col, row of layer, if it did.
func (g *Generator) edited(layer, col, row int) (int, bool) {
	k, ok := editKey(layer, col, row)
	if !ok {
		return 0, false
	}
	id, ok := g.edits[k]
	return id, ok
}

// maxLayers is how many layers an edit can tell apart.
const maxLayers = 16

// editKey is layer and cell col, row as one number, as cellKey is.
func editKey(layer, col, row int) (int, bool) {
	k, ok := cellKey(col, row)
	if !ok || layer < 0 || layer >= maxLayers {
		return 0, false
	}
	return k*maxLayers + layer, true
}
