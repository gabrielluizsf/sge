// embedwait: a method promoted from an embedded interface, whose value
// blocks only through what it calls by interface, waits as it should when
// called through another interface.
package main

import (
	"fmt"
	"sync"
)

type Getter interface{ Get() int }

// chanGetter blocks.
type chanGetter struct{ ch chan int }

func (g *chanGetter) Get() int { return <-g.ch }

// relay blocks only through the Getter it calls.
type relay struct{ inner Getter }

func (r *relay) Get() int { return r.inner.Get() + 1 }

// once embeds a Getter, as http's onceCloseListener embeds a Listener.
type once struct {
	Getter
	mu sync.Mutex
}

func main() {
	ch := make(chan int)
	go func() { ch <- 41 }()
	var g Getter = &once{Getter: &relay{inner: &chanGetter{ch}}}
	fmt.Println(g.Get())
}
