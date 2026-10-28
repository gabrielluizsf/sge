package bench

import (
	"sync"
	"testing"
)

// BenchmarkChanPingPong is one value sent each way between two
// goroutines on unbuffered channels.
func BenchmarkChanPingPong(b *testing.B) {
	ping, pong := make(chan int), make(chan int)
	go func() {
		for v := range ping {
			pong <- v
		}
	}()
	for b.Loop() {
		ping <- 1
		<-pong
	}
	close(ping)
}

// BenchmarkChanBuffered is a value through a buffered channel that never
// fills, so neither side waits.
func BenchmarkChanBuffered(b *testing.B) {
	c := make(chan int, 128)
	for b.Loop() {
		c <- 1
		<-c
	}
}

// BenchmarkWaitGroup starts and waits for eight goroutines.
func BenchmarkWaitGroup(b *testing.B) {
	for b.Loop() {
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {})
		}
		wg.Wait()
	}
}
