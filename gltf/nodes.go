package gltf

// nodes are the model's nodes, every index they hold checked, and the
// roots of its scene; a graph with a cycle, a node with two parents, or
// deeper than depth, is refused.
func nodes(doc *document, depth int) ([]Node, []int, error) {
	n := len(doc.Nodes)
	out := make([]Node, n)
	parents := make([]int, n)
	for i, dn := range doc.Nodes {
		node, err := nodeOf(doc, dn)
		if err != nil {
			return nil, nil, malformed("node %d: %v", i, err)
		}
		for _, c := range dn.Children {
			if c < 0 || c >= n || c == i {
				return nil, nil, malformed("node %d's child %d", i, c)
			}
			if parents[c]++; parents[c] > 1 {
				return nil, nil, malformed("node %d has two parents", c)
			}
		}
		out[i] = node
	}
	roots, err := sceneRoots(doc, parents)
	if err != nil {
		return nil, nil, err
	}
	seen := make([]bool, n)
	var walk func(i, d int) error
	walk = func(i, d int) error {
		if d > depth {
			return malformed("nodes nested deeper than %d", depth)
		}
		if seen[i] {
			return malformed("node %d is reached twice", i)
		}
		seen[i] = true
		for _, c := range out[i].Children {
			if err := walk(c, d+1); err != nil {
				return err
			}
		}
		return nil
	}
	for _, r := range roots {
		if err := walk(r, 0); err != nil {
			return nil, nil, err
		}
	}
	return out, roots, nil
}

// sceneRoots are the default scene's nodes, or, with no scene, every node
// with no parent.
func sceneRoots(doc *document, parents []int) ([]int, error) {
	if len(doc.Scenes) == 0 {
		var roots []int
		for i, p := range parents {
			if p == 0 {
				roots = append(roots, i)
			}
		}
		return roots, nil
	}
	s := 0
	if doc.Scene != nil {
		s = *doc.Scene
	}
	if s < 0 || s >= len(doc.Scenes) {
		return nil, malformed("scene %d of %d", s, len(doc.Scenes))
	}
	for _, r := range doc.Scenes[s].Nodes {
		if r < 0 || r >= len(parents) || parents[r] != 0 {
			return nil, malformed("scene root %d", r)
		}
	}
	return doc.Scenes[s].Nodes, nil
}

// nodeOf is a node's own fields: its mesh and skin checked, its place as
// translation, rotation and scale, from its matrix when it gives one.
func nodeOf(doc *document, dn docNode) (Node, error) {
	nd := Node{Name: dn.Name, Children: dn.Children, Mesh: -1, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}
	if dn.Mesh != nil {
		if *dn.Mesh < 0 || *dn.Mesh >= len(doc.Meshes) {
			return nd, malformed("mesh %d", *dn.Mesh)
		}
		nd.Mesh = *dn.Mesh
	}
	if dn.Skin != nil {
		if *dn.Skin < 0 || *dn.Skin >= len(doc.Skins) {
			return nd, malformed("skin %d", *dn.Skin)
		}
		nd.Skin = *dn.Skin
	}
	if len(dn.Matrix) == 16 {
		var m [16]float64
		copy(m[:], dn.Matrix)
		nd.Translation, nd.Rotation, nd.Scale = decompose(m)
		return nd, nil
	}
	if len(dn.Translation) == 3 {
		copy(nd.Translation[:], dn.Translation)
	}
	if len(dn.Rotation) == 4 {
		copy(nd.Rotation[:], dn.Rotation)
	}
	if len(dn.Scale) == 3 {
		copy(nd.Scale[:], dn.Scale)
	}
	return nd, nil
}
