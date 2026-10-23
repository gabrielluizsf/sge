package gots

import (
	"go/types"
	"strings"
)

// forward writes the class method for fn, a pointer method of a struct,
// calling its standalone function with this as the receiver: interfaces
// and TypeScript call methods on the object.
func (e *emitter) forward(w *writer, fn *types.Func, sig *types.Signature) {
	if e.wrapped(fn) {
		e.wrapper(w, fieldName(fn.Name()), "this."+e.goName(fn), sig)
	}
	params, args := []string{}, []string{"this"}
	for i := 0; i < sig.Params().Len(); i++ {
		p := "p" + itoa(i)
		params = append(params, p+": "+e.tsType(sig.Params().At(i).Type()))
		args = append(args, p)
	}
	ret := "return "
	if sig.Results().Len() == 0 {
		ret = ""
	}
	w.line("%s(%s): %s { %s%s(%s); }", e.goName(fn), strings.Join(params, ", "), e.tsResults(sig.Results()), ret, e.methodFunc(fn), strings.Join(args, ", "))
}
