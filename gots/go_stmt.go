package gots

import "go/ast"

// goStmt starts a goroutine: the function and its arguments are evaluated
// now, and the call runs once this goroutine waits or ends.
func (e *emitter) goStmt(w *writer, s *ast.GoStmt) {
	c := s.Call
	if lit, ok := ast.Unparen(c.Fun).(*ast.FuncLit); ok && len(c.Args) == 0 {
		w.line("$.go(%s);", e.funcLit(lit))
		return
	}
	w.open("{")
	if tv := e.info.Types[c.Fun]; tv.IsBuiltin() {
		w.line("$.go(() => { %s; });", e.builtin(c))
	} else {
		w.line("$.go(() => %s);", e.boundCall(w, c))
	}
	w.close("}")
}
