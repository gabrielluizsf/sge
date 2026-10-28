package gots

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

// clauseBody is a clause's statements followed, when it ends in
// fallthrough, by those of the clauses it falls into.
func clauseBody(clauses []ast.Stmt, i int) []ast.Stmt {
	var out []ast.Stmt
	for ; i < len(clauses); i++ {
		body := clauses[i].(*ast.CaseClause).Body
		if n := len(body); n > 0 {
			if b, ok := body[n-1].(*ast.BranchStmt); ok && b.Tok == token.FALLTHROUGH {
				out = append(out, body[:n-1]...)
				continue
			}
		}
		return append(out, body...)
	}
	return out
}

// chain writes an if-else chain choosing one clause: conds are the tests of
// the clauses, "" marking the default, which Go tries last wherever it is.
// The chain sits in a labelled block when a break leaves it.
func (e *emitter) chain(w *writer, label string, conds []string, body func(w *writer, i int)) {
	b := &breakable{label: label, isSwitch: true, used: label != "" && label[0] != '$'}
	if label == "" {
		b.label = e.sc.fresh("$sw")
	}
	e.fn.breaks = append(e.fn.breaks, b)
	inner := writer{indent: w.indent}
	first := true
	def := -1
	for i, c := range conds {
		if c == "" {
			def = i
			continue
		}
		if first {
			inner.open("if (%s) {", c)
			first = false
		} else {
			inner.close("} else if (%s) {", c)
			inner.indent++
		}
		body(&inner, i)
	}
	if def >= 0 {
		if first {
			inner.open("{")
		} else {
			inner.close("} else {")
			inner.indent++
		}
		body(&inner, def)
	}
	if def >= 0 || !first {
		inner.close("}")
	}
	e.fn.breaks = e.fn.breaks[:len(e.fn.breaks)-1]
	if b.used {
		w.open("%s: {", b.label)
		w.b.WriteString(indentBy(inner.String()))
		w.close("}")
		return
	}
	w.b.WriteString(inner.String())
}

func (e *emitter) switchStmt(w *writer, s *ast.SwitchStmt, label string) {
	w.open("{")
	if s.Init != nil {
		e.stmt(w, s.Init)
	}
	tag := ""
	var tagT types.Type
	if s.Tag != nil {
		tagT = e.typeOf(s.Tag)
		tag = e.sc.fresh("$tag")
		w.line("const %s = %s;", tag, e.expr(s.Tag))
	}
	clauses := s.Body.List
	conds := make([]string, len(clauses))
	for i, c := range clauses {
		cc := c.(*ast.CaseClause)
		if cc.List == nil {
			continue
		}
		cond := ""
		for j, x := range cc.List {
			var one string
			if tag == "" {
				one = e.expr(x)
			} else {
				one = e.caseEq(tag, tagT, x)
			}
			if j > 0 {
				cond += " || "
			}
			cond += one
		}
		conds[i] = cond
	}
	e.chain(w, label, conds, func(w *writer, i int) { e.stmts(w, clauseBody(clauses, i)) })
	w.close("}")
}

// caseEq tests the switch tag, bound to tag, against the case expression x.
func (e *emitter) caseEq(tag string, tagT types.Type, x ast.Expr) string {
	xt := e.typeOf(x)
	switch {
	case isNil(xt):
		return "(" + tag + " === null)"
	case isInterface(tagT):
		return "$.ifaceEq(" + tag + ", " + e.value(x, tagT) + ")"
	case isInterface(xt):
		return "$.ifaceEq(" + e.box(tag, tagT) + ", " + e.expr(x) + ")"
	}
	return e.equalOf(tag, e.expr(x), tagT)
}

// indentBy indents every line of s one step further.
func indentBy(s string) string {
	out := ""
	for _, l := range strings.SplitAfter(s, "\n") {
		if l != "" && l != "\n" {
			out += "  " + l
		} else {
			out += l
		}
	}
	return out
}
