package gots

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

var assignOps = map[token.Token]token.Token{
	token.ADD_ASSIGN: token.ADD, token.SUB_ASSIGN: token.SUB,
	token.MUL_ASSIGN: token.MUL, token.QUO_ASSIGN: token.QUO,
	token.REM_ASSIGN: token.REM, token.AND_ASSIGN: token.AND,
	token.OR_ASSIGN: token.OR, token.XOR_ASSIGN: token.XOR,
	token.SHL_ASSIGN: token.SHL, token.SHR_ASSIGN: token.SHR,
	token.AND_NOT_ASSIGN: token.AND_NOT,
}

func (e *emitter) assign(w *writer, s *ast.AssignStmt) {
	if op, ok := assignOps[s.Tok]; ok {
		e.opAssign(w, s.Lhs[0], op, e.expr(s.Rhs[0]))
		return
	}
	define := s.Tok == token.DEFINE
	if len(s.Rhs) == 1 && len(s.Lhs) > 1 {
		e.assignTuple(w, s.Lhs, s.Rhs[0], define)
		return
	}
	if len(s.Lhs) == 1 {
		e.assignOne(w, s.Lhs[0], s.Rhs[0], define)
		return
	}
	// a, b = b, a: every right-hand side is evaluated before any store.
	tmps := make([]string, len(s.Rhs))
	for i, r := range s.Rhs {
		tmps[i] = e.newTmp()
		w.line("const %s = %s;", tmps[i], e.value(r, e.typeOf(s.Lhs[i])))
	}
	for i, l := range s.Lhs {
		e.store(w, l, tmps[i], define)
	}
}

// newVar reports the variable x declares, if it is an identifier a := makes
// new rather than reuses.
func (e *emitter) newVar(x ast.Expr, define bool) *types.Var {
	id, ok := x.(*ast.Ident)
	if !ok || !define || id.Name == "_" {
		return nil
	}
	v, _ := e.info.Defs[id].(*types.Var)
	return v
}

// declareVar writes the declaration of v holding the value js.
func (e *emitter) declareVar(w *writer, v *types.Var, js string) {
	if e.predeclared[v] {
		n := e.objName(v)
		if e.boxed[v] {
			w.line("%s = new $.Cell<%s>(%s);", n, e.tsType(e.varType(v)), js)
			return
		}
		w.line("%s = %s;", n, js)
		return
	}
	n := e.declare(v)
	if e.boxed[v] {
		w.line("const %s = new $.Cell<%s>(%s);", n, e.tsType(e.varType(v)), js)
		return
	}
	w.line("let %s: %s = %s;", n, e.tsType(e.varType(v)), js)
}

// store writes the already converted value js into x.
func (e *emitter) store(w *writer, x ast.Expr, js string, define bool) {
	if v := e.newVar(x, define); v != nil {
		e.declareVar(w, v, js)
		return
	}
	lv := e.place(w, x)
	w.line("%s;", lv.set(js))
}

func (e *emitter) assignOne(w *writer, l, r ast.Expr, define bool) {
	if v := e.newVar(l, define); v != nil {
		e.declareVar(w, v, e.value(r, e.varType(v)))
		return
	}
	lv := e.place(w, l)
	if lv.blank {
		w.line("%s;", e.expr(r))
		return
	}
	if isValue(lv.t) && !lv.stores {
		// $set copies, so the value need not be copied first.
		w.line("%s;", lv.set(e.expr(r)))
		return
	}
	w.line("%s;", lv.set(e.value(r, lv.t)))
}

// assignTuple is a, b = f() and the comma-ok forms.
func (e *emitter) assignTuple(w *writer, lhs []ast.Expr, r ast.Expr, define bool) {
	var tuple string
	var types_ []types.Type
	switch r := ast.Unparen(r).(type) {
	case *ast.TypeAssertExpr:
		e.assertOkStmt(w, lhs, r, define)
		return
	case *ast.IndexExpr:
		tuple = e.mapOk(r)
		types_ = []types.Type{e.typeOf(r.X).Underlying().(*types.Map).Elem(), types.Typ[types.Bool]}
	case *ast.UnaryExpr:
		w.line("%s;", e.unsupported(r, "channel receive"))
		return
	default:
		tuple = e.expr(r)
		tup := e.typeOf(r).(*types.Tuple)
		for i := 0; i < tup.Len(); i++ {
			types_ = append(types_, tup.At(i).Type())
		}
	}
	allNew := define
	var names []string
	for _, l := range lhs {
		v := e.newVar(l, define)
		if id, ok := l.(*ast.Ident); ok && id.Name == "_" {
			names = append(names, "")
			continue
		}
		if v == nil || e.boxed[v] || e.predeclared[v] || isInterface(e.varType(v)) != isInterface(types_[len(names)]) {
			allNew = false
		}
		names = append(names, "?")
	}
	t := e.newTmp()
	w.line("const %s = %s;", t, tuple)
	if allNew {
		// let a = $t[0], b = $t[1]; not let [a, b] = f(), which V8 runs
		// through the iterator protocol.
		var ns []string
		for i, l := range lhs {
			if names[i] != "" {
				ns = append(ns, e.declare(e.newVar(l, define))+" = "+t+"["+itoa(i)+"]")
			}
		}
		if len(ns) > 0 {
			w.line("let %s;", strings.Join(ns, ", "))
		}
		return
	}
	for i, l := range lhs {
		if names[i] == "" {
			continue
		}
		to := e.typeOf(l)
		if v := e.newVar(l, define); v != nil {
			to = e.varType(v)
		}
		e.store(w, l, e.convertValue(t+"["+itoa(i)+"]", types_[i], to, false), define)
	}
}

// opAssign is x op= y, and x++ and x-- as x += 1 and x -= 1.
func (e *emitter) opAssign(w *writer, x ast.Expr, op token.Token, y string) {
	lv := e.place(w, x)
	k := defaultKind(basicKind(lv.t))
	if k == types.Uint64 && y == "1" {
		y = "1n"
	}
	if k == types.Float32 && y == "1" {
		y = "1"
	}
	var count *int
	if op == token.SHL || op == token.SHR {
		if n, ok := e.constCount(y); ok {
			count = &n
		}
	}
	plain := isString(lv.t) || k == types.Float64 || k == types.Int || k == types.Int64 || k == types.Uint || k == types.Uintptr
	if plain && (op == token.ADD || op == token.SUB || op == token.MUL && k == types.Float64 || op == token.QUO && k == types.Float64) && !strings.Contains(lv.get, "??") {
		if y == "1" && (op == token.ADD || op == token.SUB) {
			w.line("%s%s;", lv.get, map[token.Token]string{token.ADD: "++", token.SUB: "--"}[op])
			return
		}
		w.line("%s %s= %s;", lv.get, op.String(), y)
		return
	}
	w.line("%s;", lv.set(e.arith(op, lv.get, y, lv.t, count)))
}

// constCount reads a shift count written as a literal.
func (e *emitter) constCount(y string) (int, bool) {
	n := 0
	if y == "" {
		return 0, false
	}
	for _, c := range y {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}
