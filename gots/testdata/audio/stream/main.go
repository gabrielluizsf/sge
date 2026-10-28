// stream: plays half a second of a tone through antui's audio, and tells
// whether the device took about that much, and no more than it should
// have queued: the stream keeps up, and its queue stays short.
package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gabrielluizsf/antui"
)

func main() {
	pulled, phase := 0, 0.0
	a, err := antui.OpenAudio("gots audio test", 44100, func(block []float32) {
		for i := 0; i+1 < len(block); i += 2 {
			phase += 2 * math.Pi * 440 / 44100
			block[i], block[i+1] = float32(math.Sin(phase)*0.1), float32(math.Sin(phase)*0.1)
		}
		pulled += len(block) / 2
	})
	if errors.Is(err, antui.ErrNoAudio) {
		fmt.Println("no audio:", err)
		return
	}
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	time.Sleep(500 * time.Millisecond)
	a.Close()
	got := float64(pulled) / float64(a.Rate())
	fmt.Printf("rate %d, kept up %v, queue short %v\n", a.Rate(), got > 0.35, got < 0.5+0.15)
}
