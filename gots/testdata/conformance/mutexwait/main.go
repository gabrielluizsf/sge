// mutexwait: goroutines waiting their turn for locks held across a wait,
// a Mutex, an RWMutex's writers and readers, and a Cond over each.
package main

import (
	"fmt"
	"sync"
	"time"
)

type counter struct {
	sync.Mutex
	n int
}

// add holds c across a sleep, so the others wait for it.
func (c *counter) add(log chan<- string, who int) {
	c.Lock()
	defer c.Unlock()
	v := c.n
	time.Sleep(time.Millisecond)
	c.n = v + 1
	log <- fmt.Sprint("add ", who, " ", c.n)
}

func main() {
	var c counter
	log := make(chan string, 100)
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.add(log, i)
		}()
		time.Sleep(100 * time.Microsecond)
	}
	wg.Wait()
	close(log)
	n := 0
	for range log {
		n++
	}
	fmt.Println("mutex", c.n, n)

	var rw sync.RWMutex
	shared := 0
	var order []string
	var mu sync.Mutex
	note := func(s string) {
		mu.Lock()
		order = append(order, s)
		mu.Unlock()
	}
	rw.RLock()
	wg.Add(3)
	go func() {
		defer wg.Done()
		rw.Lock() // waits for the reader
		shared = 1
		note("writer")
		time.Sleep(time.Millisecond)
		rw.Unlock()
	}()
	time.Sleep(time.Millisecond)
	go func() {
		defer wg.Done()
		rw.RLock() // waits behind the waiting writer
		note(fmt.Sprint("reader sees ", shared))
		rw.RUnlock()
	}()
	go func() {
		defer wg.Done()
		var l sync.Locker = &rw
		l.Lock() // through an interface
		note("locker")
		l.Unlock()
	}()
	time.Sleep(time.Millisecond)
	note("first reader done")
	rw.RUnlock()
	wg.Wait()
	fmt.Println(order)

	var cm sync.Mutex
	cond := sync.NewCond(&cm)
	ready := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		cm.Lock()
		for !ready {
			cond.Wait()
		}
		cm.Unlock()
		fmt.Println("cond woke")
	}()
	time.Sleep(time.Millisecond)
	cm.Lock()
	ready = true
	cond.Signal()
	time.Sleep(time.Millisecond) // the waiter wakes, and waits for cm
	cm.Unlock()
	wg.Wait()
	fmt.Println("done", c.TryLock())
}
