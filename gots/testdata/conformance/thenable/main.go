// Command thenable has a struct with a field and a method named then,
// which JavaScript would take for a promise, returned from functions that
// wait and from ones that do not: they must come back as the structs they
// are, the field and the method called only when Go calls them.
package main

import "fmt"

type job struct {
	name string
	then func() string
}

type step struct{ n int }

func (s *step) then() int { return s.n * 2 }

var done = make(chan bool, 1)

// waited is a job returned by a function that waits on a channel.
func waited(name string) *job {
	done <- true
	<-done
	return &job{name: name, then: func() string { return name + " then" }}
}

// plain is a job returned by a function that does not wait.
func plain(name string) *job { return &job{name: name} }

// chooser is called through an interface, so it may wait.
type chooser interface{ choose() *step }

type pick struct{}

func (pick) choose() *step {
	done <- true
	<-done
	return &step{n: 21}
}

func main() {
	j := waited("a")
	fmt.Println(j.name, j.then())
	p := plain("b")
	fmt.Println(p.name, p.then == nil)
	var c chooser = pick{}
	s := c.choose()
	fmt.Println(s.n, s.then())
	ch := make(chan *job, 1)
	go func() { ch <- waited("c") }()
	fmt.Println((<-ch).then())
}
