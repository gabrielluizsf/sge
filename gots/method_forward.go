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
		e.wrapper(w, fieldName(fn.Name()), "this."+e.goName(fn), sig, fn)
	}
	params, args := []string{}, []string{"this"}
	for i := 0; i < sig.Params().Len(); i++ {
		p := "p" + itoa(i)
		params = append(params, p+": "+e.tsType(sig.Params().At(i).Type()))
		args = append(args, p)
	}
	ret := "return "
	if sig.Results().Len() == 0 && !e.blocks(fn) {
		ret = ""
	}
	results := e.tsResults(sig.Results())
	if e.blocks(fn) {
		results = "Promise<" + results + ">"
	}
	w.line("%s(%s): %s { %s%s(%s); }", e.goName(fn), strings.Join(params, ", "), results, ret, e.methodFunc(fn), strings.Join(args, ", "))
	if e.twin(fn) {
		w.line("%s$a(%s): Promise<%s> { return %s$a(%s); }", e.goName(fn), strings.Join(params, ", "), e.tsResults(sig.Results()), e.methodFunc(fn), strings.Join(args, ", "))
	}
}
