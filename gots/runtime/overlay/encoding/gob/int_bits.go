//go:build gots_overlay

// Overlay for the sizes of int and uint, which gob works out from uint's:
// under gots int has 64 bits, a number, and uint 32. Ints decode as int64s
// do; a uint, decoded as a uint64 is, is refused past its 32 bits as gob
// refuses it on a 32-bit machine.

package gob

import (
	"math"
	"reflect"
)

const intBits = 64

func decUint64(i *decInstr, state *decoderState, value reflect.Value) {
	v := state.decodeUint()
	if value.Kind() == reflect.Uint && math.MaxUint32 < v {
		error_(i.ovfl)
	}
	value.SetUint(v)
}
