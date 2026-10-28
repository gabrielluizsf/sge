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

// twinNatives are the hand-written functions that call what they are given
// through an interface that may block, and so have twins, name$a, as
// transpiled functions that do: fmt writing to a socket waits for it.
var twinNatives = map[string]bool{
	"fmt.Fprint":   true,
	"fmt.Fprintf":  true,
	"fmt.Fprintln": true,
	// reflect calls what it is given, which may block.
	"reflect.Value.Call":      true,
	"reflect.Value.CallSlice": true,
	// Locks wait for a goroutine that holds them across a wait.
	"sync.Mutex.Lock":    true,
	"sync.RWMutex.Lock":  true,
	"sync.RWMutex.RLock": true,
}

// twinIfaceMethods are the interface methods, as ifaceKey names them,
// that twinNatives have: a call of sync.Locker's Lock may wait.
var twinIfaceMethods = []string{"iface:Lock()", "iface:RLock()"}

// promiseTwins are the twinNatives whose twin gives back a promise only
// when it waits: a lock nobody holds is taken with no turn of the event
// loop.
var promiseTwins = map[string]bool{
	"sync.Mutex.Lock":    true,
	"sync.RWMutex.Lock":  true,
	"sync.RWMutex.RLock": true,
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
