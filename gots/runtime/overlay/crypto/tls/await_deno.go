//go:build gots_overlay

// Waiting on the promises of the gotsTls bridge, from the goroutine that
// asks: it parks on a channel the promise settles.

package tls

import (
	"errors"
	"fmt"
	"syscall/js"
)

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
		ch <- result{err: jsError(args[0])}
		return nil
	})
	p.Call("then", ok, bad)
	r := <-ch
	ok.Release()
	bad.Release()
	return r.v, r.err
}

// jsFailure is an error of node:tls, with the code Node gives it.
type jsFailure struct {
	code, message string
}

// Error quotes the message, which holds what the peer sent (names in its
// certificate), so it cannot forge lines where it is logged.
func (e *jsFailure) Error() string { return fmt.Sprintf("tls: %q", e.message) }

// jsError is the error a rejected promise gives: a Failure of the bridge,
// or whatever else JavaScript threw.
func jsError(v js.Value) error {
	if v.Type() == js.TypeObject && v.Get("failed").Truthy() {
		return &jsFailure{code: v.Get("code").String(), message: v.Get("message").String()}
	}
	if v.Type() == js.TypeObject {
		return errors.New("tls: " + fmt.Sprintf("%q", v.Get("message").String()))
	}
	return errors.New("tls: " + fmt.Sprintf("%q", v.String()))
}
