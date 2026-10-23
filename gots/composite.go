package gots

import (
	"go/ast"
	"go/constant"
	"go/types"
	"strings"
)

func (e *emitter) composite(lit *ast.CompositeLit) string {
	t := e.typeOf(lit)
	if p, ok := t.Underlying().(*types.Pointer); ok {
		t = p.Elem()
	}
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return e.structLit(lit, t, u)
	case *types.Array:
		return e.indexedLit(lit, u.Elem(), int(u.Len()))
	case *types.Slice:
		if len(lit.Elts) == 0 {
			return "$.slice(" + e.indexedLit(lit, u.Elem(), 0) + ")"
		}
		return "$.slice(" + e.indexedLit(lit, u.Elem(), -1) + ")"
	case *types.Map:
		var entries []string
		for _, el := range lit.Elts {
			kv := el.(*ast.KeyValueExpr)
			entries = append(entries, "["+e.value(kv.Key, u.Key())+", "+e.value(kv.Value, u.Elem())+"]")
		}
		if len(entries) == 0 {
			return e.newMap(u, "")
		}
		return e.newMap(u, "["+strings.Join(entries, ", ")+"]")
	}
	return e.unsupported(lit, "composite literal of this type")
}

// structLit is new C(...) with the fields in declaration order, leaving
// trailing fields to the constructor's zero defaults.
func (e *emitter) structLit(lit *ast.CompositeLit, t types.Type, st *types.Struct) string {
	vals := make([]string, st.NumFields())
	for i, el := range lit.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			name := kv.Key.(*ast.Ident).Name
			for j := 0; j < st.NumFields(); j++ {
				if st.Field(j).Name() == name {
					vals[j] = e.value(kv.Value, st.Field(j).Type())
				}
			}
			continue
		}
		vals[i] = e.value(el, st.Field(i).Type())
	}
	n := len(vals)
	for n > 0 && vals[n-1] == "" {
		n--
	}
	for i := 0; i < n; i++ {
		if vals[i] == "" {
			vals[i] = e.zero(st.Field(i).Type())
		}
	}
	return "new " + e.structClass(t) + "(" + strings.Join(vals[:n], ", ") + ")"
}

// indexedLit is the backing array of an array or slice literal of n
// elements, n < 0 meaning as many as the literal reaches. Elements may name
// their index.
func (e *emitter) indexedLit(lit *ast.CompositeLit, elem types.Type, n int) string {
	type entry struct {
		i int
		v string
	}
	var entries []entry
	next, max := 0, 0
	for _, el := range lit.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			k, _ := constant.Int64Val(e.info.Types[kv.Key].Value)
			next = int(k)
			el = kv.Value
		}
		entries = append(entries, entry{next, e.value(el, elem)})
		next++
		if next > max {
			max = next
		}
	}
	if n < 0 {
		n = max
	}
	dense := len(entries) == n
	for i, en := range entries {
		if en.i != i {
			dense = false
		}
	}
	if dense {
		vals := make([]string, n)
		for i, en := range entries {
			vals[i] = en.v
		}
		return e.arrayLit(elem, vals)
	}
	var zero string
	if b := backing(elem); b != "" {
		zero = "new " + b + "(" + itoa(n) + ")"
	} else {
		zero = "$.fillArray(" + itoa(n) + ", () => " + e.zero(elem) + ")"
	}
	if len(entries) == 0 {
		return zero
	}
	var sets []string
	for _, en := range entries {
		sets = append(sets, "a["+itoa(en.i)+"] = "+en.v+";")
	}
	return "((a: any) => { " + strings.Join(sets, " ") + " return a; })(" + zero + ")"
}
