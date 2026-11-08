package lineage

// Entry is one person of the tree as plain data, with whether they died.
type Entry struct {
	Person
	Dead bool
}

// Entries are everyone in the tree, in the order they were added, each
// after their parents.
func (t *Tree) Entries() []Entry {
	es := make([]Entry, 0, len(t.order))
	for _, id := range t.order {
		m := t.people[id]
		es = append(es, Entry{Person: m.Person, Dead: m.dead})
	}
	return es
}

// Restore is the tree of entries, as Entries gave them; an error if one is
// not a valid addition, as Add says.
func Restore(entries []Entry) (*Tree, error) {
	t := &Tree{}
	for _, e := range entries {
		if err := t.Add(e.Person); err != nil {
			return nil, err
		}
		if e.Dead {
			t.Die(e.ID)
		}
	}
	return t, nil
}
