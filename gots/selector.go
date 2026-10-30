package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

func (e *emitter) selector(x *ast.SelectorExpr) string {
	s := e.info.Selections[x]
	if s == nil {
		// A qualified identifier: pkg.Name.
		return e.ident(x.Sel)
	}
	switch s.Kind() {
	case types.FieldVal:
		js, _ := e.fieldPath(e.expr(x.X), e.T(s.Recv()), s.Index())
		return js
	case types.MethodVal:
		return e.methodValue(x, s)
	case types.MethodExpr:
		return e.methodExpr(s)
	}
	return e.unsupported(x, "selector")
}

// fieldPath follows the field indices path from base, of type t, through
// embedded fields; pointers on the way need no step, as a pointer to a
// struct is the struct.
func (e *emitter) fieldPath(base string, t types.Type, path []int) (string, types.Type) {
	for _, i := range path {
		st := structOf(t)
		base += "." + e.fieldJS(st, i)
		t = st.Field(i).Type()
	}
	return base, t
}

// methodFunc is the standalone function implementing a method: of a named
// type that is not a struct, which has no class to hold it, or a pointer
// method of a struct, which Go code calls without its class.
func (e *emitter) methodFunc(fn *types.Func) string {
	recv := types.Unalias(derefType(e.recvType(fn))).(*types.Named)
	if isStruct(recv) {
		return e.namedRef(recv) + "$$" + fn.Name()
	}
	return e.namedRef(recv) + "$" + fn.Name()
}

func namedOf(t types.Type) *types.Named {
	n, _ := types.Unalias(t).(*types.Named)
	if n != nil && n.Origin() != nil {
		return n.Origin()
	}
	return n
}

// onClass reports whether calls to fn are method calls on an object: its
// receiver is a struct or an interface.
func onClass(fn *types.Func) bool {
	recv := derefType(fn.Type().(*types.Signature).Recv().Type())
	return isStruct(recv) || isInterface(recv)
}

// methodValue is x.M used as a function value. The receiver is evaluated,
// and copied for a value method, when the value is made.
func (e *emitter) methodValue(x *ast.SelectorExpr, s *types.Selection) string {
	fn := s.Obj().(*types.Func)
	recv := e.receiver(x, s)
	// bound is the function value calling call(r, ...a), with a twin
	// calling twin(r, ...a) when there is one.
	bound := func(call, twin string) string {
		f := "(...a: any[]) => " + call
		if twin != "" {
			f = "$.twins(" + f + ", async (...a: any[]) => " + twin + ")"
		}
		return "((r) => " + f + ")(" + recv + ")"
	}
	if onClass(fn) && !e.freeMethod(s) {
		if !isPointer(e.recvType(fn)) && !isInterface(e.recvType(fn)) {
			recv = e.cloneOf(recv, derefType(e.typeOf(x.X)))
		}
		if isInterface(e.recvType(fn)) {
			sig := fn.Type().(*types.Signature)
			twin := ""
			if e.ifaceMayBlock(fn) {
				twin = fromIface("(await $.callA(r, "+jsQuote(fieldName(fn.Name()))+", ...a))", sig)
			}
			return bound(fromIface("r."+fieldName(fn.Name())+"(...a)", sig), twin)
		}
		twin := ""
		if e.twin(fn) {
			twin = "r." + e.goName(fn) + "$a(...a)"
		}
		return bound("r."+e.goName(fn)+"(...a)", twin)
	}
	if r := e.recvType(fn); !isPointer(r) && isValue(r) {
		recv = e.cloneOf(recv, r)
	}
	twin := ""
	if e.twin(fn) {
		twin = e.methodFunc(fn) + "$a(r, ...a)"
	}
	return bound(e.methodFunc(fn)+"(r, ...a)", twin)
}

// methodExpr is T.M or (*T).M: a function taking the receiver first.
func (e *emitter) methodExpr(s *types.Selection) string {
	fn := s.Obj().(*types.Func)
	if len(s.Index()) > 1 || onClass(fn) {
		path := ""
		if len(s.Index()) > 1 {
			p, _ := e.fieldPath("", e.T(s.Recv()), s.Index()[:len(s.Index())-1])
			path = p
		}
		return "((r: any, ...a: any[]) => r" + path + "." + e.goName(fn) + "(...a))"
	}
	return e.methodFunc(fn)
}

// receiver is the receiver argument of a call to the method selected by x:
// for a method on a class, the object to call it on; otherwise the value or
// pointer the standalone method function takes first.
func (e *emitter) receiver(x *ast.SelectorExpr, s *types.Selection) string {
	fn := s.Obj().(*types.Func)
	path := s.Index()[:len(s.Index())-1]
	if e.handWritten(e.T(s.Recv())) {
		// The TypeScript standard library has promoted methods on the
		// outer type itself.
		path = nil
	}
	base, t := e.fieldPath(e.expr(x.X), e.T(s.Recv()), path)
	if onClass(fn) {
		return base
	}
	_, ptrRecv := e.recvType(fn).(*types.Pointer)
	switch {
	case ptrRecv && isPointer(t):
		return base
	case ptrRecv:
		if len(path) == 0 {
			return e.addr(x.X)
		}
		dot := strings.LastIndex(base, ".")
		return "new $.FieldPtr(" + base[:dot] + ", " + jsQuote(base[dot+1:]) + ")"
	case isPointer(t):
		return derefPtr(base, t)
	}
	return base
}

func isPointer(t types.Type) bool {
	_, ok := t.Underlying().(*types.Pointer)
	return ok
}

func jsQuote(s string) string { return `"` + s + `"` }

// handWritten reports whether t is a type of the TypeScript standard
// library rather than of a transpiled package.
func (e *emitter) handWritten(t types.Type) bool {
	n := namedOf(derefType(t))
	if n == nil || n.Obj().Pkg() == nil {
		return false
	}
	p := e.prog.Packages[n.Obj().Pkg().Path()]
	return p != nil && p.Std && !p.Transpile
}
