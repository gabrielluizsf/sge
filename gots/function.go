package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// function writes prefix(params): R { body }, or (params): R => { body } for
// a closure. prologue are statements to run before the body, such as
// binding the receiver. fresh starts a new set of local names; a closure
// shares its enclosing function's.
func (e *emitter) function(w *writer, prefix string, sig *types.Signature, ft *ast.FuncType, body *ast.BlockStmt, prologue []string, arrow, fresh bool) {
	outerSc, outerFn := e.sc, e.fn
	if fresh || e.sc == nil {
		e.sc = newScope(e.pkgNames)
	}
	e.fn = &funcState{sig: sig, ft: ft, body: body, parent: outerFn}
	defer func() { e.sc, e.fn = outerSc, outerFn }()

	var params, boxes []string
	for _, field := range ft.Params.List {
		t := e.typeOf(field.Type)
		names := field.Names
		if len(names) == 0 {
			names = []*ast.Ident{nil}
		}
		for _, id := range names {
			var obj types.Object
			if id != nil {
				obj = e.info.Defs[id]
			}
			ts := e.tsType(t)
			if field.Type != nil {
				if _, variadic := field.Type.(*ast.Ellipsis); variadic {
					ts = e.tsType(types.NewSlice(t))
				}
			}
			if obj == nil || obj.Name() == "_" {
				params = append(params, e.sc.fresh("_")+": "+ts)
				continue
			}
			n := e.declare(obj)
			if e.boxed[obj] {
				params = append(params, n+"$in: "+ts)
				boxes = append(boxes, "const "+n+" = new $.Cell("+n+"$in);")
				continue
			}
			params = append(params, n+": "+ts)
		}
	}
	ret := e.tsResults(sig.Results())
	head := "(" + strings.Join(params, ", ") + "): " + ret
	if arrow {
		w.open("%s => {", head)
	} else {
		w.open("%s%s {", prefix, head)
	}
	for _, p := range prologue {
		w.line("%s", p)
	}
	for _, b := range boxes {
		w.line("%s", b)
	}
	if ft.Results != nil {
		e.namedResults(w, ft.Results)
	}
	if body == nil {
		e.errorf(ft.Pos(), "%s has no body: it needs an overlay", strings.TrimPrefix(prefix, "export function "))
		w.line(`$.panic("gots: %s has no body");`, prefix)
	} else if defers(body) {
		e.deferringBody(w, body)
	} else {
		e.stmts(w, body.List)
	}
	w.close("}")
}

// namedResults declares a function's named results, at their zero values.
func (e *emitter) namedResults(w *writer, results *ast.FieldList) {
	for _, field := range results.List {
		for _, id := range field.Names {
			obj := e.info.Defs[id]
			var n string
			if obj == nil || id.Name == "_" {
				n = e.sc.fresh("r")
			} else {
				n = e.declare(obj)
			}
			t := e.typeOf(field.Type)
			if obj != nil && e.boxed[obj] {
				w.line("const %s = new $.Cell(%s);", n, e.zero(t))
				n += ".v"
			} else {
				w.line("let %s: %s = %s;", n, e.tsType(t), e.zero(t))
			}
			e.fn.results = append(e.fn.results, n)
		}
	}
}

// defers reports whether body has a defer statement of its own, outside any
// closure.
func defers(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.DeferStmt:
			found = true
		}
		return !found
	})
	return found
}

// deferringBody runs the body under a Frame: deferred calls run however it
// ends, a recovered panic returns the results as they stand.
func (e *emitter) deferringBody(w *writer, body *ast.BlockStmt) {
	e.fn.frame = e.sc.fresh("$f")
	w.line("const %s = new $.Frame();", e.fn.frame)
	w.open("try {")
	e.stmts(w, body.List)
	w.close("} catch ($e) {")
	w.indent++
	w.line("%s.fail($e);", e.fn.frame)
	w.close("} finally {")
	w.indent++
	w.line("%s.unwind();", e.fn.frame)
	results := e.fn.sig.Results()
	switch {
	case len(e.fn.results) > 0:
		w.line("return %s;", e.resultTuple(e.fn.results))
	case results.Len() > 0:
		var zeros []string
		for i := 0; i < results.Len(); i++ {
			zeros = append(zeros, e.zero(results.At(i).Type()))
		}
		w.line("if (%s.recovered) return %s;", e.fn.frame, e.resultTuple(zeros))
	}
	w.close("}")
}

func (e *emitter) resultTuple(vals []string) string {
	if len(vals) == 1 {
		return vals[0]
	}
	return "[" + strings.Join(vals, ", ") + "]"
}
