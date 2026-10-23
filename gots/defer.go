package gots

import (
	"go/ast"
	"go/types"
)

// deferStmt registers a call on the function's Frame. The function and its
// arguments are evaluated now, as Go does; only the call waits.
func (e *emitter) deferStmt(w *writer, s *ast.DeferStmt) {
	c := s.Call
	if lit, ok := ast.Unparen(c.Fun).(*ast.FuncLit); ok && len(c.Args) == 0 {
		w.line("%s.defer(%s);", e.fn.frame, e.funcLit(lit))
		return
	}
	w.open("{")
	var args []string
	if tv := e.info.Types[c.Fun]; tv.IsBuiltin() {
		w.line("%s.defer(() => { %s; });", e.fn.frame, e.builtin(c))
		w.close("}")
		return
	}
	sig := e.typeOf(c.Fun).Underlying().(*types.Signature)
	for _, a := range e.args(c, sig) {
		t := e.newTmp()
		w.line("const %s = %s;", t, a)
		args = append(args, t)
	}
	call := e.callWith(c, args)
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
	w.line("%s.defer(() => { %s; });", e.fn.frame, call)
	w.close("}")
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
