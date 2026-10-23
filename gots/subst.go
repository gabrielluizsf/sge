package gots

import "go/types"

// T is t with the type parameters of the generic code being emitted
// replaced by the instance's type arguments. Each instance of a generic
// function or type is emitted on its own, so emitted code only ever sees
// concrete types.
func (e *emitter) T(t types.Type) types.Type {
	if t == nil || len(e.subst) == 0 {
		return t
	}
	return subst(t, e.subst, e.ctxt)
}

// varType is the type of a variable, after substitution.
func (e *emitter) varType(v types.Object) types.Type { return e.T(v.Type()) }

// recvType is the receiver type of the method fn, after substitution.
func (e *emitter) recvType(fn *types.Func) types.Type {
	return e.T(fn.Type().(*types.Signature).Recv().Type())
}

func subst(t types.Type, m map[*types.TypeParam]types.Type, ctxt *types.Context) types.Type {
	s := func(t types.Type) types.Type { return subst(t, m, ctxt) }
	switch t := t.(type) {
	case *types.TypeParam:
		if r, ok := m[t]; ok {
			return r
		}
		return t
	case *types.Basic:
		return t
	case *types.Alias:
		return s(types.Unalias(t))
	case *types.Named:
		args := t.TypeArgs()
		if args.Len() == 0 {
			return t
		}
		changed := false
		next := make([]types.Type, args.Len())
		for i := 0; i < args.Len(); i++ {
			next[i] = s(args.At(i))
			changed = changed || next[i] != args.At(i)
		}
		if !changed {
			return t
		}
		inst, err := types.Instantiate(ctxt, t.Origin(), next, false)
		if err != nil {
			return t
		}
		return inst
	case *types.Pointer:
		return types.NewPointer(s(t.Elem()))
	case *types.Slice:
		return types.NewSlice(s(t.Elem()))
	case *types.Array:
		return types.NewArray(s(t.Elem()), t.Len())
	case *types.Map:
		return types.NewMap(s(t.Key()), s(t.Elem()))
	case *types.Chan:
		return types.NewChan(t.Dir(), s(t.Elem()))
	case *types.Tuple:
		return substTuple(t, m, ctxt)
	case *types.Signature:
		var recv *types.Var
		if r := t.Recv(); r != nil {
			recv = types.NewParam(r.Pos(), r.Pkg(), r.Name(), s(r.Type()))
		}
		return types.NewSignatureType(recv, nil, nil, substTuple(t.Params(), m, ctxt), substTuple(t.Results(), m, ctxt), t.Variadic())
	case *types.Struct:
		fields := make([]*types.Var, t.NumFields())
		tags := make([]string, t.NumFields())
		for i := range fields {
			f := t.Field(i)
			fields[i] = types.NewField(f.Pos(), f.Pkg(), f.Name(), s(f.Type()), f.Embedded())
			tags[i] = t.Tag(i)
		}
		return types.NewStruct(fields, tags)
	case *types.Interface:
		var methods []*types.Func
		for i := 0; i < t.NumExplicitMethods(); i++ {
			fn := t.ExplicitMethod(i)
			methods = append(methods, types.NewFunc(fn.Pos(), fn.Pkg(), fn.Name(), s(fn.Type()).(*types.Signature)))
		}
		var embeds []types.Type
		for i := 0; i < t.NumEmbeddeds(); i++ {
			embeds = append(embeds, s(t.EmbeddedType(i)))
		}
		return types.NewInterfaceType(methods, embeds).Complete()
	case *types.Union:
		terms := make([]*types.Term, t.Len())
		for i := range terms {
			terms[i] = types.NewTerm(t.Term(i).Tilde(), s(t.Term(i).Type()))
		}
		return types.NewUnion(terms)
	}
	return t
}

func substTuple(t *types.Tuple, m map[*types.TypeParam]types.Type, ctxt *types.Context) *types.Tuple {
	if t == nil {
		return nil
	}
	vars := make([]*types.Var, t.Len())
	for i := range vars {
		v := t.At(i)
		vars[i] = types.NewParam(v.Pos(), v.Pkg(), v.Name(), subst(v.Type(), m, ctxt))
	}
	return types.NewTuple(vars...)
}
