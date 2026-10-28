package main

import "fmt"

type Path struct{ pts []int }

func (p *Path) Length() int {
	if p == nil {
		return 0
	}
	return len(p.pts)
}

type Tree struct {
	Left, Right *Tree
	Val         int
}

func (t *Tree) Sum() int {
	if t == nil {
		return 0
	}
	return t.Val + t.Left.Sum() + t.Right.Sum()
}

type Holder struct{ *Path }

func main() {
	var none *Path
	fmt.Println(none.Length(), (&Path{pts: []int{1, 2}}).Length())
	t := &Tree{Val: 1, Left: &Tree{Val: 2}, Right: &Tree{Val: 3, Left: &Tree{Val: 4}}}
	fmt.Println(t.Sum())
	var h Holder
	fmt.Println(h.Length())
	f := none.Length
	fmt.Println(f())
}

type loc struct{ name string }

var utc = &loc{"UTC"}

// get swaps a nil receiver for a default, reassigning the receiver itself.
func (l *loc) get() *loc {
	if l == nil {
		return utc
	}
	return l
}

func (l *loc) String() string {
	l = l.get()
	return l.name
}

func init() {
	var none *loc
	fmt.Println(none.String(), (&loc{"Local"}).String())
}
