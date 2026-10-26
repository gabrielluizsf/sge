//go:build gots_overlay

// Overlay for the recover every gob decode and encode runs under. Data
// nested deep in a recursive type the sender defines, a list of a few
// thousand links in some 6 KB, is walked by recursion no depth limit of
// gob's holds: Go's stack grows, a JavaScript one overflows, and the
// program ended. The overflow, which gots recovers as a runtime error, is
// the stream's fault, and is reported as gob reports a bad stream; by the
// time it gets here the stack has unwound.

package gob

import (
	"errors"
	"runtime"
	"strings"
)

var errTooDeep = errors.New("gob: invalid nesting depth")

func catchError(err *error) {
	if e := recover(); e != nil {
		ge, ok := e.(gobError)
		if !ok {
			if re, isRT := e.(runtime.Error); isRT && strings.Contains(re.Error(), "Maximum call stack size exceeded") {
				*err = errTooDeep
				return
			}
			panic(e)
		}
		*err = ge.err
	}
}
