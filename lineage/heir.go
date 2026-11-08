package lineage

// Heir is who takes over from id when they die: their eldest living
// child, or, none living, the eldest living heir of their children in
// turn; false when there is no one.
func (t *Tree) Heir(id string) (string, bool) {
	m := t.people[id]
	if m == nil {
		return "", false
	}
	best, found := "", false
	for _, c := range m.children {
		if !t.people[c].dead && (!found || t.people[c].Born < t.people[best].Born) {
			best, found = c, true
		}
	}
	if found {
		return best, true
	}
	for _, c := range m.children {
		if h, ok := t.Heir(c); ok {
			return h, true
		}
	}
	return "", false
}
