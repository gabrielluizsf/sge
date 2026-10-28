// Goroutines and channels: every program here prints the same however the
// goroutines are scheduled.
package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

func pingPong() {
	ping, pong := make(chan int), make(chan int)
	go func() {
		for v := range ping {
			pong <- v * 2
		}
		close(pong)
	}()
	for i := 1; i <= 3; i++ {
		ping <- i
		fmt.Println("pong", <-pong)
	}
	close(ping)
	_, ok := <-pong
	fmt.Println("pong closed:", !ok)
}

func buffered() {
	c := make(chan string, 3)
	c <- "a"
	c <- "b"
	fmt.Println(len(c), cap(c))
	close(c)
	for s := range c {
		fmt.Print(s, " ")
	}
	v, ok := <-c
	fmt.Printf("%q %v\n", v, ok)
}

func workers() {
	jobs := make(chan int)
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- j * j
			}
		}()
	}
	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	var got []int
	for r := range results {
		got = append(got, r)
	}
	sort.Ints(got)
	fmt.Println("squares", got)
}

func selects() {
	a, b := make(chan int), make(chan string)
	done := make(chan struct{})
	go func() {
		a <- 1
		b <- "two"
		close(done)
	}()
	var seen []string
	for len(seen) < 2 {
		select {
		case n := <-a:
			seen = append(seen, fmt.Sprint("a", n))
		case s, ok := <-b:
			seen = append(seen, fmt.Sprint("b", s, ok))
		}
	}
	<-done
	fmt.Println(seen)

	// A select with a default does not wait.
	empty := make(chan int)
	select {
	case v := <-empty:
		fmt.Println("got", v)
	default:
		fmt.Println("nothing ready")
	}

	// A nil channel is never ready.
	var never chan int
	full := make(chan int, 1)
	select {
	case <-never:
		fmt.Println("nil ready?")
	case full <- 5:
		fmt.Println("sent", <-full)
	}
}

func timeouts() {
	slow := make(chan int)
	select {
	case v := <-slow:
		fmt.Println("got", v)
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timed out")
	}
	start := time.Now()
	time.Sleep(15 * time.Millisecond)
	fmt.Println("slept enough:", time.Since(start) >= 15*time.Millisecond)

	t := time.NewTicker(5 * time.Millisecond)
	n := 0
	for range t.C {
		if n++; n == 3 {
			t.Stop()
			break
		}
	}
	fmt.Println("ticks", n)
}

var mu sync.Mutex

func counter() {
	total := 0
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Go(func() {
			mu.Lock()
			total += i
			mu.Unlock()
		})
	}
	wg.Wait()
	fmt.Println("total", total)
}

// generator hands out values from a goroutine.
func generator(words []string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for _, w := range words {
			c <- strings.ToUpper(w)
		}
	}()
	return c
}

type stage struct{ name string }

func (s *stage) run(in <-chan int, out chan<- int) {
	for v := range in {
		out <- v + 1
	}
	close(out)
}

func pipeline() {
	c0, c1, c2 := make(chan int), make(chan int), make(chan int)
	go (&stage{"one"}).run(c0, c1)
	go (&stage{"two"}).run(c1, c2)
	go func() {
		for i := 0; i < 3; i++ {
			c0 <- i * 10
		}
		close(c0)
	}()
	for v := range c2 {
		fmt.Print(v, " ")
	}
	fmt.Println()
}

// recovered panics inside a goroutine, and a deferred call that waits.
func recovered() {
	errs := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errs <- fmt.Errorf("recovered: %v", r)
			}
		}()
		var m map[string]int
		m["x"] = 1
	}()
	fmt.Println(<-errs != nil)

	out := make(chan string, 1)
	func() {
		defer func() { out <- "deferred send" }()
	}()
	fmt.Println(<-out)
}

func cond() {
	var m sync.Mutex
	c := sync.NewCond(&m)
	ready := false
	done := make(chan bool)
	go func() {
		m.Lock()
		for !ready {
			c.Wait()
		}
		m.Unlock()
		done <- true
	}()
	m.Lock()
	ready = true
	c.Broadcast()
	m.Unlock()
	fmt.Println("cond", <-done)
}

// fanIn reads two producers until both are done, through a function value
// that blocks.
func fanIn() {
	produce := func(n int) <-chan int {
		c := make(chan int)
		go func() {
			for i := 0; i < n; i++ {
				c <- i
			}
			close(c)
		}()
		return c
	}
	a, b := produce(3), produce(2)
	sum := 0
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil
				continue
			}
			sum += v
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			sum += v
		}
	}
	fmt.Println("fan-in", sum)
}

var errClosed = errors.New("closed")

func sendOnClosed() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	c := make(chan int, 1)
	close(c)
	c <- 1
	return errClosed
}

// closureCalls calls a closure that blocks through the variable that
// holds it, which makes the caller wait too.
func closureCalls() int {
	c := make(chan int, 1)
	put := func(v int) { c <- v }
	take := func() int { return <-c }
	put(7)
	return take() * 6
}

// polling loops on a select with a default until another goroutine says
// stop: Go's scheduler takes turns with it, and so must this.
func polling() {
	done := make(chan struct{})
	go func() {
		close(done)
	}()
	spins := 0
	for {
		select {
		case <-done:
			fmt.Println("polling stopped:", spins > 0 || spins == 0)
			return
		default:
			spins++
		}
	}
}

func main() {
	polling()
	fmt.Println("closure", closureCalls())
	pingPong()
	buffered()
	workers()
	selects()
	timeouts()
	counter()
	for w := range generator([]string{"go", "to", "ts"}) {
		fmt.Print(w, " ")
	}
	fmt.Println()
	pipeline()
	recovered()
	cond()
	fanIn()
	fmt.Println(sendOnClosed())
	go func() {
		time.Sleep(time.Hour)
		fmt.Println("never printed: main returned")
	}()
}
