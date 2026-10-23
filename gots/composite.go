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
	if e.handWritten(t) {
		// A hand-written class has its own fields: set the named ones.
		var sets []string
		for _, el := range lit.Elts {
			kv := el.(*ast.KeyValueExpr)
			name := kv.Key.(*ast.Ident).Name
			for j := 0; j < st.NumFields(); j++ {
				if st.Field(j).Name() == name {
					sets = append(sets, fieldName(name)+": "+e.value(kv.Value, st.Field(j).Type()))
				}
			}
		}
		return "$Object.assign(new " + e.structClass(t) + "(), { " + strings.Join(sets, ", ") + " })"
	}
	vals := make([]string, st.NumFields())
	var promoted []string
	for i, el := range lit.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			name := kv.Key.(*ast.Ident).Name
			direct := false
			for j := 0; j < st.NumFields(); j++ {
				if st.Field(j).Name() == name {
					vals[j] = e.value(kv.Value, st.Field(j).Type())
					direct = true
				}
			}
			if !direct {
				promoted = append(promoted, e.promotedSet("o", t, name, kv.Value))
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
	made := "new " + e.structClass(t) + "(" + strings.Join(vals[:n], ", ") + ")"
	if len(promoted) == 0 {
		return made
	}
	return "((o: any) => { " + strings.Join(promoted, " ") + " return o; })(" + made + ")"
}

// promotedSet is the statement setting the field name that the struct
// type t gets from an embedded field, on the object o: the literal form
// Go 1.27 allows. An embedded pointer on the way that is nil is given a
// value to hold the field.
func (e *emitter) promotedSet(o string, t types.Type, name string, v ast.Expr) string {
	obj, index, _ := types.LookupFieldOrMethod(t, false, e.tpkg, name)
	f, ok := obj.(*types.Var)
	if !ok {
		return ""
	}
	var steps []string
	path := o
	cur := t
	for _, i := range index[:len(index)-1] {
		fs := structOf(cur).Field(i)
		path += "." + e.fieldJS(structOf(cur), i)
		if p, isPtr := fs.Type().Underlying().(*types.Pointer); isPtr {
			steps = append(steps, "if ("+path+" === null) "+path+" = "+e.zero(p.Elem())+";")
		}
		cur = fs.Type()
	}
	path += "." + e.fieldJS(structOf(cur), index[len(index)-1])
	steps = append(steps, e.setOf(path, e.value(v, f.Type()), f.Type())+";")
	return strings.Join(steps, " ")
}

// indexedLit is the backing array of an array or slice literal of n
// elements, n < 0 meaning as many as the literal reaches. Elements may name
// their index.
func (e *emitter) indexedLit(lit *ast.CompositeLit, elem types.Type, n int) string {
	return e.indexedLitIn(lit, elem, n, backing(elem))
}

// indexedLitIn is indexedLit held in a b typed array, or a plain array if
// b is "".
func (e *emitter) indexedLitIn(lit *ast.CompositeLit, elem types.Type, n int, b string) string {
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
		if b != "" {
			return b + ".of(" + strings.Join(vals, ", ") + ")"
		}
		return "[" + strings.Join(vals, ", ") + "]"
	}
	var zero string
	if b != "" {
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
