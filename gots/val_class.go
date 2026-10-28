package gots

import "go/types"

// valClassOf is the class an interface holds values of t in, as its
// descriptor's valCls names it: the type's $Val class when it has value
// methods, $.Val otherwise. It is "" where only the runtime knows, for
// pointers and hand-written types.
func (e *emitter) valClassOf(t types.Type) string {
	if isPointer(t) || isInterface(t) || e.handWritten(t) {
		return ""
	}
	n, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return "$.Val"
	}
	if types.NewMethodSet(n).Len() == 0 {
		return "$.Val"
	}
	return e.namedRef(n) + "$Val"
}
