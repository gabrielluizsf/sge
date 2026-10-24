package gots

import "go/types"

// blockingNatives are the hand-written functions that block, by
// nativeName.
var blockingNatives = map[string]bool{
	"sync.WaitGroup.Wait": true,
	"sync.Cond.Wait":      true,
	"runtime.Gosched":     true,
	"testing.T.Run":       true,
}

// nativeName is fn as pkg.Func or pkg.Type.Method.
func nativeName(fn *types.Func) string {
	if fn.Pkg() == nil {
		return ""
	}
	sig := fn.Type().(*types.Signature)
	if r := sig.Recv(); r != nil {
		if n := namedOf(derefType(r.Type())); n != nil {
			return fn.Pkg().Path() + "." + n.Obj().Name() + "." + fn.Name()
		}
		return ""
	}
	return fn.Pkg().Path() + "." + fn.Name()
}
