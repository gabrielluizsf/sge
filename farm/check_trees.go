package farm

import "fmt"

// checkStumps makes the stumps: each on the cell of the tree it was, which
// has to be standing there now, none twice.
func (c *checker) checkStumps() error {
	if len(c.s.Stumps) > maxStumps {
		return fmt.Errorf("save: %d stumps", len(c.s.Stumps))
	}
	seen := map[int]bool{}
	for _, s := range c.s.Stumps {
		if err := firstError(checkCell("stump", s.Col, s.Row), finite("stump time", s.Left, 0, regrowTime(s.Tree))); err != nil {
			return err
		}
		k := cellKey(s.Col, s.Row)
		if c.oldWorld() && c.f.thingAt(s.Col, s.Row) != s.Tree {
			continue // ground a town has since built on
		}
		if !isTree(s.Tree) && !isRock(s.Tree) || c.f.thingAt(s.Col, s.Row) != s.Tree || seen[k] {
			return fmt.Errorf("save: no tree %d to have been felled at %d, %d", s.Tree, s.Col, s.Row)
		}
		seen[k] = true
		c.out.stumps = append(c.out.stumps, &stump{col: s.Col, row: s.Row, tree: s.Tree, left: s.Left})
	}
	return nil
}

// checkHits makes the strokes of the trees being cut: each on a tree
// standing, a number the axes would not have felled it at.
func (c *checker) checkHits() error {
	if len(c.s.Hits) > maxStumps {
		return fmt.Errorf("save: %d trees being cut", len(c.s.Hits))
	}
	c.out.hits = map[int]int{}
	felled := map[int]bool{}
	for _, s := range c.out.stumps {
		felled[cellKey(s.col, s.row)] = true
	}
	for _, h := range c.s.Hits {
		if err := firstError(checkCell("tree", h.Col, h.Row), inRange("strokes", h.N, 1, 3)); err != nil {
			return err
		}
		k := cellKey(h.Col, h.Row)
		if c.oldWorld() && !isTree(c.f.thingAt(h.Col, h.Row)) && !isRock(c.f.thingAt(h.Col, h.Row)) {
			continue // ground a town has since built on
		}
		_, twice := c.out.hits[k]
		if !isTree(c.f.thingAt(h.Col, h.Row)) && !isRock(c.f.thingAt(h.Col, h.Row)) || twice || felled[k] {
			return fmt.Errorf("save: no tree or rock being cut at %d, %d", h.Col, h.Row)
		}
		c.out.hits[k] = h.N
	}
	return nil
}
