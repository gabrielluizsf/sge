package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// bindRange gives a range variable its value for this iteration: a new
// variable for :=, a store for =.
func (e *emitter) bindRange(w *writer, x ast.Expr, js string, from types.Type, tok token.Token) {
	if x == nil {
		return
	}
	if id, ok := x.(*ast.Ident); ok && id.Name == "_" {
		return
	}
	if tok == token.DEFINE {
		v := e.newVar(x, true)
		e.declareVar(w, v, e.convertValue(js, from, e.varType(v), false))
		return
	}
	lv := e.place(w, x)
	w.line("%s;", lv.set(e.convertValue(js, from, lv.t, false)))
}

func (e *emitter) rangeStmt(w *writer, s *ast.RangeStmt, label string) {
	t := e.typeOf(s.X)
	head := labelPrefix(label)
	i, n := e.sc.fresh("$i"), e.sc.fresh("$n")
	switch u := t.Underlying().(type) {
	case *types.Basic:
		if u.Info()&types.IsString != 0 {
			e.rangeString(w, s, head)
			return
		}
		// range n
		zero, one := "0", ""
		if isBig(t) {
			zero = "0n"
		}
		_ = one
		e.loop(w, head+"for (let "+i+" = "+zero+", "+n+" = "+e.expr(s.X)+"; "+i+" < "+n+"; "+i+"++)", func() {
			e.bindRange(w, s.Key, i, t, s.Tok)
			e.stmts(w, s.Body.List)
		})
	case *types.Slice:
		sl := e.sc.fresh("$s")
		e.loop(w, head+"for (let "+sl+" = "+e.expr(s.X)+", "+i+" = 0, "+n+" = $.len("+sl+"); "+i+" < "+n+"; "+i+"++)", func() {
			e.bindRange(w, s.Key, i, types.Typ[types.Int], s.Tok)
			e.bindRange(w, s.Value, e.cloneOf(sl+".a["+sl+".o + "+i+"]", u.Elem()), u.Elem(), s.Tok)
			e.stmts(w, s.Body.List)
		})
	case *types.Array, *types.Pointer:
		var elem types.Type
		if a, ok := u.(*types.Array); ok {
			elem = a.Elem()
		} else {
			elem = u.(*types.Pointer).Elem().Underlying().(*types.Array).Elem()
		}
		a := e.sc.fresh("$a")
		e.loop(w, head+"for (let "+a+" = "+e.expr(s.X)+", "+i+" = 0; "+i+" < "+a+".length; "+i+"++)", func() {
			e.bindRange(w, s.Key, i, types.Typ[types.Int], s.Tok)
			e.bindRange(w, s.Value, e.cloneOf(a+"["+i+"]", elem), elem, s.Tok)
			e.stmts(w, s.Body.List)
		})
	case *types.Map:
		k, v := e.sc.fresh("$k"), e.sc.fresh("$v")
		e.loop(w, head+"for (const ["+k+", "+v+"] of $.entries("+e.expr(s.X)+"))", func() {
			e.bindRange(w, s.Key, k, u.Key(), s.Tok)
			e.bindRange(w, s.Value, e.cloneOf(v, u.Elem()), u.Elem(), s.Tok)
			e.stmts(w, s.Body.List)
		})
	case *types.Signature:
		e.rangeFunc(w, s, u, head)
	default:
		w.line("%s;", e.unsupported(s, "range over this type"))
	}
}

// rangeString walks a string rune by rune, as Go does.
func (e *emitter) rangeString(w *writer, s *ast.RangeStmt, head string) {
	str, i, wd, r := e.sc.fresh("$str"), e.sc.fresh("$i"), e.sc.fresh("$w"), e.sc.fresh("$r")
	e.loop(w, head+"for (let "+str+" = "+e.expr(s.X)+", "+i+" = 0, "+wd+" = 0; "+i+" < "+str+".length; "+i+" += "+wd+")", func() {
		w.line("let %s = %s.charCodeAt(%s);", r, str, i)
		w.line("%s = 1;", wd)
		w.line("if (%s >= 128) [%s, %s] = $.str.decodeRune(%s, %s);", r, r, wd, str, i)
		e.bindRange(w, s.Key, i, types.Typ[types.Int], s.Tok)
		e.bindRange(w, s.Value, r, types.Typ[types.Int32], s.Tok)
		e.stmts(w, s.Body.List)
	})
}

// rangeFunc is range over an iterator function: the body becomes the yield
// callback, returning false to stop.
func (e *emitter) rangeFunc(w *writer, s *ast.RangeStmt, sig *types.Signature, head string) {
	if head != "" || hasBranchOut(s.Body) {
		w.line("%s;", e.unsupported(s, "break, continue or return inside range over a function"))
		return
	}
	params := sig.Params().At(0).Type().(*types.Signature).Params()
	vars := []string{}
	for i := 0; i < params.Len(); i++ {
		vars = append(vars, e.sc.fresh("$y"))
	}
	list := ""
	for i, v := range vars {
		if i > 0 {
			list += ", "
		}
		list += v + ": any"
	}
	w.open("%s((%s) => {", e.expr(s.X), list)
	keys := []ast.Expr{s.Key, s.Value}
	for i, v := range vars {
		e.bindRange(w, keys[i], v, params.At(i).Type(), s.Tok)
	}
	e.stmts(w, s.Body.List)
	w.line("return true;")
	w.close("});")
}

// hasBranchOut reports whether body leaves the loop it is in by any way
// other than finishing an iteration.
func hasBranchOut(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			found = true
		case *ast.BranchStmt:
			found = true
		case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			_ = n
		}
		return !found
	})
	return found
}
