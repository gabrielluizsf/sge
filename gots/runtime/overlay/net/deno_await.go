//go:build gots_overlay

// Waiting on the promises of the gotsNet bridge, from the goroutine that
// asks: it parks on a channel the promise settles.

package net

import "syscall/js"

// awaitJS is what the promise p gives, or the error it is rejected with.
func awaitJS(p js.Value) (js.Value, error) {
	type result struct {
		v   js.Value
		err error
	}
	ch := make(chan result, 1)
	var ok, bad js.Func
	ok = js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- result{v: args[0]}
		return nil
	})
	bad = js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- result{err: denoError(args[0])}
		return nil
	})
	p.Call("then", ok, bad)
	r := <-ch
	ok.Release()
	bad.Release()
	return r.v, r.err
}
