// Package inventory is what someone carries: a bag of items by name, so
// many of each, that one can put in, count, and take out.
package inventory

// Bag holds items by name. Its zero value is an empty bag.
type Bag struct {
	counts map[string]int
	// names are the items ever put in, in the order first put in, so going
	// through a bag goes the same way every run.
	names []string
}

// Stack is n of an item.
type Stack struct {
	Item string
	N    int
}

// List is stacks of items: a wish, a recipe, a bag's content.
type List []Stack
