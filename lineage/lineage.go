// Package lineage is a family tree in a game: who is whose child, who is
// kin to whom, so kin never court, and who is heir when one dies, so a
// child takes over from their parent.
package lineage

import (
	"errors"
	"fmt"
)

// Person is one of the tree: their id, their parents', "" for one not in
// the tree, and the day they were born.
type Person struct {
	ID      string
	Parents [2]string
	Born    int
}

// Tree is a family tree, its people by id, each added after their parents.
type Tree struct {
	people map[string]*member
}

// member is a person of the tree, their children in the order born, and
// whether they have died.
type member struct {
	Person
	children []string
	dead     bool
}

// ErrUnknown is an id not in the tree.
var ErrUnknown = errors.New("lineage: no such person")

// Add puts p in the tree: their parents, if any, must be in it already,
// and their id new, so the tree has no loops.
func (t *Tree) Add(p Person) error {
	if t.people == nil {
		t.people = map[string]*member{}
	}
	if p.ID == "" || t.people[p.ID] != nil {
		return fmt.Errorf("lineage: %q is already in the tree", p.ID)
	}
	for _, id := range p.Parents {
		if id != "" && t.people[id] == nil {
			return fmt.Errorf("%w: parent %q", ErrUnknown, id)
		}
	}
	if p.Parents[0] != "" && p.Parents[0] == p.Parents[1] {
		return fmt.Errorf("lineage: %q twice a parent", p.Parents[0])
	}
	t.people[p.ID] = &member{Person: p}
	for _, id := range p.Parents {
		if id != "" {
			t.people[id].children = append(t.people[id].children, p.ID)
		}
	}
	return nil
}

// Die marks id as dead.
func (t *Tree) Die(id string) error {
	m := t.people[id]
	if m == nil {
		return ErrUnknown
	}
	m.dead = true
	return nil
}

// Children are id's children, eldest first.
func (t *Tree) Children(id string) []string {
	m := t.people[id]
	if m == nil {
		return nil
	}
	return append([]string(nil), m.children...)
}
