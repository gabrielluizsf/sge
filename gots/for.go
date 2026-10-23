package gots

import (
	"go/ast"
	"go/token"
	"strings"
)

func labelPrefix(label string) string {
	if label == "" {
		return ""
	}
	return label + ": "
}

// loop writes a loop's body with the loop on the break stack. head may
// start with the loop's Go label; a loop without one is given a label
// only when a break or continue inside has to name it, so the body is
// written first and the head after.
func (e *emitter) loop(w *writer, head string, body func()) {
	b := &breakable{}
	if i := strings.Index(head, ": "); i > 0 && !strings.ContainsAny(head[:i], " (") {
		b.label = head[:i]
		b.used = true
		head = head[i+2:]
	} else {
		b.label = e.sc.fresh("$loop")
	}
	e.fn.breaks = append(e.fn.breaks, b)
	before := w.b.String()
	w.b.Reset()
	w.indent++
	body()
	text := w.b.String()
	w.indent--
	w.b.Reset()
	w.b.WriteString(before)
	e.fn.breaks = e.fn.breaks[:len(e.fn.breaks)-1]
	if b.used {
		w.line("%s: %s {", b.label, head)
	} else {
		w.line("%s {", head)
	}
	w.b.WriteString(text)
	w.line("}")
}

func (e *emitter) forStmt(w *writer, s *ast.ForStmt, label string) {
	init := ""
	wrapped := false
	if s.Init != nil {
		if h, ok := e.forInit(s.Init); ok {
			init = h
		} else {
			wrapped = true
			w.open("{")
			e.stmt(w, s.Init)
		}
	}
	cond := ""
	if s.Cond != nil {
		cond = e.expr(s.Cond)
	}
	post := ""
	if s.Post != nil {
		post = e.postExpr(s.Post)
	}
	e.loop(w, labelPrefix(label)+"for ("+init+"; "+cond+"; "+post+")", func() { e.stmts(w, s.Body.List) })
	if wrapped {
		w.close("}")
	}
}

// forInit is a loop's i := 0 as a let in the loop's head, which gives each
// iteration its own i, as Go 1.22 does.
func (e *emitter) forInit(s ast.Stmt) (string, bool) {
	as, ok := s.(*ast.AssignStmt)
	if !ok || as.Tok != token.DEFINE || len(as.Lhs) != len(as.Rhs) {
		return "", false
	}
	var parts []string
	for i, l := range as.Lhs {
		v := e.newVar(l, true)
		if v == nil || e.boxed[v] {
			return "", false
		}
		parts = append(parts, e.declare(v)+": "+e.tsType(e.varType(v))+" = "+e.value(as.Rhs[i], e.varType(v)))
	}
	return "let " + strings.Join(parts, ", "), true
}

// postExpr is a loop's post statement as the expression a for head needs.
func (e *emitter) postExpr(s ast.Stmt) string {
	if as, ok := s.(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) > 1 && len(as.Lhs) == len(as.Rhs) {
		// a, b = b, a as a destructuring assignment, which evaluates every
		// right-hand side before storing.
		var pre writer
		var gets, vals []string
		for i, l := range as.Lhs {
			lv := e.place(&pre, l)
			if lv.blank || isValue(lv.t) {
				return e.unsupported(s, "this post statement")
			}
			gets = append(gets, lv.get)
			vals = append(vals, e.value(as.Rhs[i], lv.t))
		}
		if pre.String() != "" {
			return e.unsupported(s, "post statement with effects")
		}
		return "[" + strings.Join(gets, ", ") + "] = [" + strings.Join(vals, ", ") + "]"
	}
	var out writer
	e.stmt(&out, s)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 {
		return e.unsupported(s, "post statement with effects")
	}
	return strings.TrimSuffix(lines[0], ";")
}
