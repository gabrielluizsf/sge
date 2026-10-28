package gots

import (
	"go/types"
	"strconv"
	"strings"
)

// rtString is t written the way Go's reflect writes it, which %T and %v
// print: interface {} with a space, any spelled out, packages by name.
func rtString(t types.Type) string {
	switch t := t.(type) {
	case *types.Alias:
		return rtString(types.Unalias(t))
	case *types.Named:
		s := t.Obj().Name()
		if t.Obj().Pkg() != nil {
			s = t.Obj().Pkg().Name() + "." + s
		}
		if args := t.TypeArgs(); args.Len() > 0 {
			var as []string
			for i := 0; i < args.Len(); i++ {
				as = append(as, rtString(args.At(i)))
			}
			s += "[" + strings.Join(as, ",") + "]"
		}
		return s
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			return "unsafe.Pointer"
		}
		return types.Typ[defaultKind(t.Kind())].Name()
	case *types.Pointer:
		return "*" + rtString(t.Elem())
	case *types.Slice:
		return "[]" + rtString(t.Elem())
	case *types.Array:
		return "[" + strconv.FormatInt(t.Len(), 10) + "]" + rtString(t.Elem())
	case *types.Map:
		return "map[" + rtString(t.Key()) + "]" + rtString(t.Elem())
	case *types.Chan:
		switch t.Dir() {
		case types.SendOnly:
			return "chan<- " + rtString(t.Elem())
		case types.RecvOnly:
			return "<-chan " + rtString(t.Elem())
		}
		return "chan " + rtString(t.Elem())
	case *types.Signature:
		return "func" + rtSignature(t)
	case *types.Interface:
		if t.NumMethods() == 0 {
			return "interface {}"
		}
		var ms []string
		for i := 0; i < t.NumMethods(); i++ {
			m := t.Method(i)
			ms = append(ms, m.Name()+rtSignature(m.Type().(*types.Signature)))
		}
		return "interface { " + strings.Join(ms, "; ") + " }"
	case *types.Struct:
		if t.NumFields() == 0 {
			return "struct {}"
		}
		var fs []string
		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			if f.Embedded() {
				fs = append(fs, rtString(f.Type()))
			} else {
				fs = append(fs, f.Name()+" "+rtString(f.Type()))
			}
		}
		return "struct { " + strings.Join(fs, "; ") + " }"
	case *types.TypeParam:
		return t.Obj().Name()
	}
	return t.String()
}

func rtSignature(sig *types.Signature) string {
	var ps []string
	for i := 0; i < sig.Params().Len(); i++ {
		p := rtString(sig.Params().At(i).Type())
		if sig.Variadic() && i == sig.Params().Len()-1 {
			p = "..." + strings.TrimPrefix(p, "[]")
		}
		ps = append(ps, p)
	}
	s := "(" + strings.Join(ps, ", ") + ")"
	switch r := sig.Results(); r.Len() {
	case 0:
	case 1:
		s += " " + rtString(r.At(0).Type())
	default:
		var rs []string
		for i := 0; i < r.Len(); i++ {
			rs = append(rs, rtString(r.At(i).Type()))
		}
		s += " (" + strings.Join(rs, ", ") + ")"
	}
	return s
}
