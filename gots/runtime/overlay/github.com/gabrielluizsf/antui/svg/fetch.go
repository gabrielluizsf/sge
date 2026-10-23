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

// pictureTimeout is antui's pictureHTTP timeout, in milliseconds.
const pictureTimeout = 10000

// fetchPictureAsync reads target through the runtime's gotsFetchBytes,
// which holds getPicture's limits: a 200 answer only, at most
// maxPictureBytes read, and the timeout.
func fetchPictureAsync(target string, e *pictureBytes) {
	var onBody, onError js.Func
	onError = js.FuncOf(func(this js.Value, args []js.Value) any {
		finishPicture(target, e, nil, fmt.Errorf("%s", args[0].Get("message").String()))
		return nil
	})
	onBody = js.FuncOf(func(this js.Value, args []js.Value) any {
		data := make([]byte, args[0].Length())
		js.CopyBytesToGo(data, args[0])
		finishPicture(target, e, data, nil)
		return nil
	})
	js.Global().Call("gotsFetchBytes", target, maxPictureBytes, pictureTimeout).Call("then", onBody, onError)
}
