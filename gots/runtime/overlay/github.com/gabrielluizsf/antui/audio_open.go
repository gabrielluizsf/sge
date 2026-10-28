//go:build gots_overlay

// Overlay for audio_other.go: under gots sound goes out through the
// gotsAudio bridge, SDL2's queue on Deno and Web Audio in a browser,
// reached with syscall/js. Without the bridge, or a device, OpenAudio
// fails with ErrNoAudio, and a game carries on in silence.

package antui

import (
	"fmt"
	"syscall/js"
)

func openAudio(name string, a *Audio) (driver, error) {
	bridge := js.Global().Get("gotsAudio")
	if bridge.IsUndefined() {
		return nil, fmt.Errorf("%w: there is no audio bridge", ErrNoAudio)
	}
	dev := bridge.Call("open", name, a.rate)
	if dev.Type() == js.TypeString {
		return nil, fmt.Errorf("%w: %s", ErrNoAudio, dev.String())
	}
	a.rate = dev.Get("rate").Int()
	s := &jsAudio{dev: dev, a: a}
	go s.run()
	return s, nil
}
