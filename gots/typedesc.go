package gots

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"
)

// desc is an expression for t's run-time descriptor. Descriptors of unnamed
// types are hoisted into module-level constants so building one costs
// nothing at run time.
func (e *emitter) desc(t types.Type) string {
	d := e.descInline(t)
	if strings.HasPrefix(d, "$.") && !strings.Contains(d, "(") {
		return d
	}
	if _, named := t.(*types.Named); named {
		return d
	}
	if c, ok := e.typeConsts[d]; ok {
		return c
	}
	c := "$T" + strconv.Itoa(len(e.typeConsts)+1)
	e.typeConsts[d] = c
	e.typeOrder = append(e.typeOrder, d)
	return c
}

var basicDescs = map[types.BasicKind]string{
	types.Bool: "$.Bool", types.Int: "$.Int", types.Int8: "$.Int8",
	types.Int16: "$.Int16", types.Int32: "$.Int32", types.Int64: "$.Int64",
	types.Uint: "$.Uint", types.Uint8: "$.Uint8", types.Uint16: "$.Uint16",
	types.Uint32: "$.Uint32", types.Uint64: "$.Uint64", types.Uintptr: "$.Uintptr",
	types.Float32: "$.Float32", types.Float64: "$.Float64",
	types.Complex64: "$.Complex64", types.Complex128: "$.Complex128",
	types.String: "$.String", types.UnsafePointer: "$.UnsafePointer",
}

// descInline is desc without hoisting, for descriptors built while the
// module loads, before hoisted constants exist.
func (e *emitter) descInline(t types.Type) string {
	switch t := t.(type) {
	case *types.Basic:
		if t.Kind() == types.UntypedNil {
			return "null"
		}
		return basicDescs[defaultKind(t.Kind())]
	case *types.Alias:
		return e.descInline(types.Unalias(t))
	case *types.Named:
		if t.Obj().Pkg() == nil {
			if t.Obj().Name() == "error" {
				return `$.ifaceType("error", ["Error"])`
			}
			return `$.ifaceType("comparable", [])`
		}
		if t.TypeArgs().Len() > 0 && isInterface(t) {
			u := t.Underlying().(*types.Interface)
			return fmt.Sprintf("$.ifaceType(%q, %s, %s)", typeString(t), methodList(u), e.methodSigs(u))
		}
		name := e.namedRef(t)
		if isStruct(t) {
			return name + ".$type"
		}
		return name + "$type"
	case *types.Pointer:
		return "$.ptrTo(" + e.descInline(t.Elem()) + ")"
	case *types.Slice:
		return "$.sliceOf(" + e.descInline(t.Elem()) + ")"
	case *types.Array:
		return fmt.Sprintf("$.arrayOf(%s, %d)", e.descInline(t.Elem()), t.Len())
	case *types.Map:
		return fmt.Sprintf("$.mapOf(%s, %s)", e.descInline(t.Key()), e.descInline(t.Elem()))
	case *types.Chan:
		dir := "chan "
		switch t.Dir() {
		case types.SendOnly:
			dir = "chan<- "
		case types.RecvOnly:
			dir = "<-chan "
		}
		return fmt.Sprintf("$.chanOf(%s, %q)", e.descInline(t.Elem()), dir)
	case *types.Signature:
		var ins, outs []string
		for i := 0; i < t.Params().Len(); i++ {
			ins = append(ins, e.descInline(t.Params().At(i).Type()))
		}
		for i := 0; i < t.Results().Len(); i++ {
			outs = append(outs, e.descInline(t.Results().At(i).Type()))
		}
		return fmt.Sprintf("$.funcType(%q, () => [%s], () => [%s], %v)", typeString(t), strings.Join(ins, ", "), strings.Join(outs, ", "), t.Variadic())
	case *types.Interface:
		return fmt.Sprintf("$.ifaceType(%q, %s, %s)", typeString(t), methodList(t), e.methodSigs(t))
	case *types.Struct:
		return e.anonClass(t) + ".$type"
	case *types.TypeParam:
		return "$.Int"
	case *types.Tuple:
		return "null"
	}
	return "null"
}

// methodSigs is a thunk of the descriptors of an interface's method
// types, in methodList's order, for reflect.
func (e *emitter) methodSigs(t *types.Interface) string {
	var sigs []string
	for i := 0; i < t.NumMethods(); i++ {
		sigs = append(sigs, e.descInline(e.T(t.Method(i).Type())))
	}
	return "() => [" + strings.Join(sigs, ", ") + "]"
}

// methodList is the names an interface asks for, as a JavaScript array.
func methodList(t *types.Interface) string {
	var names []string
	for i := 0; i < t.NumMethods(); i++ {
		names = append(names, strconv.Quote(fieldName(t.Method(i).Name())))
	}
	return "[" + strings.Join(names, ", ") + "]"
}
