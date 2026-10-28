package gots

import (
	"go/ast"
	"go/types"
)

// deferStmt registers a call on the function's Frame. The function and its
// arguments are evaluated now, as Go does; only the call waits. In a
// function that may block the deferred call hands back what it returns,
// a promise for one that blocks too.
func (e *emitter) deferStmt(w *writer, s *ast.DeferStmt) {
	c := s.Call
	if lit, ok := ast.Unparen(c.Fun).(*ast.FuncLit); ok && len(c.Args) == 0 {
		w.line("%s.defer(%s);", e.fn.frame, e.funcLit(lit))
		return
	}
	w.open("{")
	if tv := e.info.Types[c.Fun]; tv.IsBuiltin() {
		w.line("%s.defer(() => { %s; });", e.fn.frame, e.builtin(c))
		w.close("}")
		return
	}
	call := e.boundCall(w, c)
	if e.fn.async {
		w.line("%s.defer(() => %s);", e.fn.frame, call)
	} else {
		w.line("%s.defer(() => { %s; });", e.fn.frame, call)
	}
	w.close("}")
}

// boundCall is the call c with its function, receiver and arguments
// evaluated now, into constants written to w, for a defer or go statement
// to make later. Arguments are copies, as the call is not made now.
func (e *emitter) boundCall(w *writer, c *ast.CallExpr) string {
	sig := e.typeOf(c.Fun).Underlying().(*types.Signature)
	var args []string
	e.deferring = true
	bound := e.args(c, sig)
	e.deferring = false
	for _, a := range bound {
		t := e.newTmp()
		w.line("const %s = %s;", t, a)
		args = append(args, t)
	}
	e.bindingCall = true
	call := e.callWith(c, args)
	e.bindingCall = false
	if sel, ok := ast.Unparen(c.Fun).(*ast.SelectorExpr); ok && !e.simple(sel.X) {
		if s := e.info.Selections[sel]; s != nil && s.Kind() == types.MethodVal {
			// Bind the receiver now too.
			r := e.newTmp()
			w.line("const %s = %s;", r, e.methodValue(sel, s))
			call = r + "(" + joinArgs(args) + ")"
		}
	} else if _, isIdent := ast.Unparen(c.Fun).(*ast.Ident); !isIdent && !ok {
		f := e.newTmp()
		w.line("const %s = %s;", f, e.expr(c.Fun))
		call = f + "(" + joinArgs(args) + ")"
	}
	return call
}

func joinArgs(args []string) string {
	out := ""
	for i, a := range args {
		if i > 0 {
			out += ", "
		}
		out += a
	}
	return out
}
