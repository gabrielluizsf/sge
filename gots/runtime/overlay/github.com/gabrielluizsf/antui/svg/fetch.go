//go:build gots_overlay

// Overlay for the pictures svg reads from the network: instead of a
// goroutine blocking on net/http, an http(s) address is read with
// JavaScript's fetch, and the drawings waiting on it are told when its
// promise settles. Everything else is read where it is asked for.

package svg

import (
	"fmt"
	"syscall/js"
)

var pictureHTTP = struct{}{}

// getPicture is only reached through fetchPicture for an address the
// overlay already sent to fetch.
func getPicture(target string) ([]byte, error) {
	return nil, fmt.Errorf("%s is read asynchronously", target)
}

func pictureBytesFor(target string) *pictureBytes {
	pictureByteCache.mu.Lock()
	if e := pictureByteCache.entries[target]; e != nil {
		pictureByteCache.mu.Unlock()
		return e
	}
	e := &pictureBytes{}
	pictureByteCache.entries[target] = e
	pictureByteCache.mu.Unlock()
	if s := hrefScheme(target); s == "http" || s == "https" {
		fetchPictureAsync(target, e)
	} else {
		fetchPictureLater(target, e)
	}
	return e
}

func fetchPictureLater(target string, e *pictureBytes) {
	data, err := fetchPicture(target)
	finishPicture(target, e, data, err)
}

// finishPicture decodes what was read and tells every drawing waiting.
func finishPicture(target string, e *pictureBytes, data []byte, fetchErr error) {
	read := pictureRead{data: data, err: fetchErr}
	pic, err := read.decoded(nil, target)

	e.mu.Lock()
	e.read, e.done = read, true
	waiting := e.waiting
	e.waiting = nil
	e.mu.Unlock()

	if fetchErr != nil {
		forgetPictureBytes(target)
	}
	for _, tell := range waiting {
		tell(pic, err)
	}
}

func fetchPictureAsync(target string, e *pictureBytes) {
	var onBody, onAnswer, onError js.Func
	onError = js.FuncOf(func(this js.Value, args []js.Value) any {
		finishPicture(target, e, nil, fmt.Errorf("%s", args[0].Get("message").String()))
		return nil
	})
	onBody = js.FuncOf(func(this js.Value, args []js.Value) any {
		buf := js.Global().Get("Uint8Array").New(args[0])
		if buf.Length() > maxPictureBytes {
			finishPicture(target, e, nil, fmt.Errorf("it is bigger than %d bytes", maxPictureBytes))
			return nil
		}
		data := make([]byte, buf.Length())
		js.CopyBytesToGo(data, buf)
		finishPicture(target, e, data, nil)
		return nil
	})
	onAnswer = js.FuncOf(func(this js.Value, args []js.Value) any {
		resp := args[0]
		if !resp.Get("ok").Bool() {
			finishPicture(target, e, nil, fmt.Errorf("it answered %d %s", resp.Get("status").Int(), resp.Get("statusText").String()))
			return nil
		}
		resp.Call("arrayBuffer").Call("then", onBody, onError)
		return nil
	})
	js.Global().Call("fetch", target).Call("then", onAnswer, onError)
}
