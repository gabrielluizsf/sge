//go:build gots_overlay

package antui

import (
	"syscall/js"
	"time"
)

// jsAudio is a stream the gotsAudio bridge plays. The bridge only queues
// what it is given, so a goroutine keeps lead's worth queued, a block of
// tick's at a time: no more, so the game is heard soon after it plays a
// sound, and the queue never grows.
type jsAudio struct {
	dev js.Value
	a   *Audio
}

const (
	audioLead = 60 * time.Millisecond
	audioTick = 10 * time.Millisecond
)

func (s *jsAudio) run() {
	block := make([]float32, 2*s.a.rate*int(audioTick)/int(time.Second))
	lead := audioLead.Seconds()
	for !s.a.stopped() {
		for s.dev.Call("queued").Float() < lead && !s.a.stopped() {
			s.a.pull(block)
			s.dev.Call("push", block)
		}
		time.Sleep(audioTick)
	}
}

func (s *jsAudio) close() error {
	s.dev.Call("close")
	return nil
}

func (s *jsAudio) latency() time.Duration { return audioLead }
