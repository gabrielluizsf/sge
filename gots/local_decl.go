package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// localDecl is a var, const or type declaration inside a function.
func (e *emitter) localDecl(w *writer, d *ast.GenDecl) {
	switch d.Tok {
	case token.CONST, token.IMPORT:
		// Constants are written out where they are used.
	case token.VAR:
		for _, spec := range d.Specs {
			vs := spec.(*ast.ValueSpec)
			e.varSpec(w, vs)
		}
	case token.TYPE:
		for _, spec := range d.Specs {
			e.localType(w, spec.(*ast.TypeSpec))
		}
	}
}

func (e *emitter) varSpec(w *writer, vs *ast.ValueSpec) {
	switch {
	case len(vs.Values) == 0:
		for _, id := range vs.Names {
			if id.Name == "_" {
				continue
			}
			v := e.info.Defs[id].(*types.Var)
			e.declareVar(w, v, e.zero(e.varType(v)))
		}
	case len(vs.Values) == 1 && len(vs.Names) > 1:
		lhs := make([]ast.Expr, len(vs.Names))
		for i, id := range vs.Names {
			lhs[i] = id
		}
		e.assignTuple(w, lhs, vs.Values[0], true)
	default:
		for i, id := range vs.Names {
			if id.Name == "_" {
				w.line("%s;", e.expr(vs.Values[i]))
				continue
			}
			v := e.info.Defs[id].(*types.Var)
			e.declareVar(w, v, e.value(vs.Values[i], e.varType(v)))
		}
	}
}

// localType declares a type inside a function: a class for a struct, a
// descriptor for anything else.
func (e *emitter) localType(w *writer, ts *ast.TypeSpec) {
	tn, ok := e.info.Defs[ts.Name].(*types.TypeName)
	if !ok || tn.IsAlias() {
		return
	}
	n := e.sc.fresh(tn.Name())
	e.localTypes[tn] = n
	if st, ok := tn.Type().Underlying().(*types.Struct); ok {
		e.structClassDecl(w, n, st, false, nil)
		w.line("%s.$type = $.structType(%q, true, %s, null, () => %s);", n, typeString(tn.Type()), n, e.fieldDescs(st))
		return
	}
	w.line("const %s$type = $.named(%q, %s, null);", n, typeString(tn.Type()), e.descInline(tn.Type().Underlying()))
}
