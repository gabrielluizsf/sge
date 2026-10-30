package lineage

// Kin is whether a and b are family too near to court: one is the
// other, or their parent or grandparent, or they share a parent or a
// grandparent, as siblings, aunts, nephews and first cousins do.
func (t *Tree) Kin(a, b string) bool {
	if a == b {
		return true
	}
	up := func(id string) map[string]bool { // id, and their parents and grandparents
		near := map[string]bool{id: true}
		for _, p := range t.parents(id) {
			near[p] = true
			for _, g := range t.parents(p) {
				near[g] = true
			}
		}
		return near
	}
	na, nb := up(a), up(b)
	if na[b] || nb[a] {
		return true
	}
	for id := range na {
		if id != a && nb[id] {
			return true
		}
	}
	return false
}

// parents are id's parents in the tree.
func (t *Tree) parents(id string) []string {
	m := t.people[id]
	if m == nil {
		return nil
	}
	var out []string
	for _, p := range m.Parents {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
