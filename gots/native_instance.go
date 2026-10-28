package gots

import (
	"go/types"
	"strings"
)

// typedNatives are the hand-written generic functions that are given the
// descriptors of their type arguments; the others have them erased.
var typedNatives = map[string]bool{
	"reflect.TypeFor":    true,
	"reflect.TypeAssert": true,
	"errors.AsType":      true,
}

// nativeInstance is a hand-written generic function instantiated with
// args: the function with the descriptors of its type arguments bound
// ahead of its arguments, reflect.TypeFor[T]() being TypeFor(T's).
func (e *emitter) nativeInstance(fn *types.Func, args *types.TypeList) string {
	descs := make([]string, args.Len())
	for i := range descs {
		descs[i] = e.desc(e.T(args.At(i)))
	}
	return e.objName(fn) + ".bind(null, " + strings.Join(descs, ", ") + ")"
}
