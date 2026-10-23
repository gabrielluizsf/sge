package gots

import (
	"go/ast"
	"go/types"
)

func (e *emitter) returnStmt(w *writer, s *ast.ReturnStmt) {
	results := e.fn.sig.Results()
	if len(s.Results) == 0 {
		if len(e.fn.results) > 0 {
			w.line("return %s;", e.resultTuple(e.fn.results))
		} else {
			w.line("return;")
		}
		return
	}
	var vals []string
	if len(s.Results) == 1 && results.Len() > 1 {
		// return f(), passing on a call's results.
		tup := e.typeOf(s.Results[0]).(*types.Tuple)
		t := e.newTmp()
		w.line("const %s = %s;", t, e.expr(s.Results[0]))
		for i := 0; i < results.Len(); i++ {
			vals = append(vals, e.convertValue(t+"["+itoa(i)+"]", tup.At(i).Type(), results.At(i).Type(), false))
		}
	} else {
		for i, r := range s.Results {
			vals = append(vals, e.value(r, results.At(i).Type()))
		}
	}
	if e.fn.frame != "" && len(e.fn.results) > 0 {
		// Deferred calls may still change named results, so the values go
		// into them and the function returns whatever they end up as.
		for i, v := range vals {
			w.line("%s = %s;", e.fn.results[i], v)
		}
		w.line("return %s;", e.resultTuple(e.fn.results))
		return
	}
	w.line("return %s;", e.resultTuple(vals))
}
