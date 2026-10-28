package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// instance is one instantiation of a generic function or type, emitted in
// the module that uses it.
type instance struct {
	obj   types.Object // the generic function or type name
	targs []types.Type
	name  string
}

// generic reports whether obj is a generic function or type that this
// module has to instantiate itself: a hand-written standard library
// package provides its generics already.
func (e *emitter) generic(obj types.Object) bool {
	if obj.Pkg() != nil {
		if p := e.prog.Packages[obj.Pkg().Path()]; p != nil && p.Std && hasStd(p.Path) {
			return false
		}
	}
	switch o := obj.(type) {
	case *types.Func:
		return o.Type().(*types.Signature).TypeParams().Len() > 0
	case *types.TypeName:
		n, ok := o.Type().(*types.Named)
		return ok && n.TypeParams().Len() > 0
	}
	return false
}

// instName is the binding of obj instantiated with targs, queueing the
// instance to be emitted with this module.
func (e *emitter) instName(obj types.Object, targs *types.TypeList) string {
	args := make([]types.Type, targs.Len())
	var strs []string
	for i := range args {
		args[i] = e.T(targs.At(i))
		strs = append(strs, rtString(args[i]))
	}
	key := obj.Pkg().Path() + "." + obj.Name() + "[" + strings.Join(strs, ",") + "]"
	if in, ok := e.insts[key]; ok {
		return in.name
	}
	base := mangle(obj.Name()) + "$$" + sanitize(strings.Join(strs, "_"))
	if obj.Pkg() != e.tpkg {
		base = sanitize(obj.Pkg().Name()) + "$" + base
	}
	name := base
	for i := 2; e.pkgNames[name]; i++ {
		name = base + "$" + itoa(i)
	}
	e.pkgNames[name] = true
	in := &instance{obj: obj, targs: args, name: name}
	e.insts[key] = in
	e.pendingInsts = append(e.pendingInsts, in)
	return name
}

// sanitize turns a type string into identifier characters.
func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		case r == '*':
			b.WriteString("P")
		case r == '[':
			b.WriteString("S")
		case r == ']', r == ' ':
		default:
			b.WriteString("_")
		}
	}
	return b.String()
}

// namedRef is how emitted code names the type n: its class or binding, or
// its instance's for a generic type.
func (e *emitter) namedRef(n *types.Named) string {
	if n.TypeArgs().Len() > 0 && e.generic(n.Origin().Obj()) {
		return e.instName(n.Origin().Obj(), n.TypeArgs())
	}
	return e.objName(n.Obj())
}

// declOf finds the declaration of a generic function or type, and the
// type information it was checked with.
func (e *emitter) declOf(obj types.Object) (ast.Node, *types.Info) {
	look := func(files []*ast.File, info *types.Info) ast.Node {
		for _, f := range files {
			for _, d := range f.Decls {
				switch d := d.(type) {
				case *ast.FuncDecl:
					if info.Defs[d.Name] == obj {
						return d
					}
				case *ast.GenDecl:
					for _, s := range d.Specs {
						if ts, ok := s.(*ast.TypeSpec); ok && info.Defs[ts.Name] == obj {
							return ts
						}
					}
				}
			}
		}
		return nil
	}
	if n := look(e.files, e.info0); n != nil {
		return n, e.info0
	}
	if p := e.prog.Packages[obj.Pkg().Path()]; p != nil {
		if n := look(p.Files, p.Info); n != nil {
			return n, p.Info
		}
	}
	return nil, nil
}

// emitInstances writes every queued instance, including those the
// instances themselves need.
func (e *emitter) emitInstances(td *typeDecls, funcs *writer) {
	for len(e.pendingInsts) > 0 {
		in := e.pendingInsts[0]
		e.pendingInsts = e.pendingInsts[1:]
		node, info := e.declOf(in.obj)
		if node == nil {
			e.errs = append(e.errs, errorf("no declaration for generic %s", in.obj.Name()))
			continue
		}
		if !e.boxedDone[info] {
			e.boxedDone[info] = true
			e.findBoxedIn(e.filesOf(in.obj.Pkg()), info)
		}
		saved, savedSubst, savedFn, savedSc := e.info, e.subst, e.fn, e.sc
		e.info = info
		e.fn, e.sc = nil, nil
		switch d := node.(type) {
		case *ast.FuncDecl:
			fn := in.obj.(*types.Func)
			e.subst = bind(fn.Type().(*types.Signature).TypeParams(), in.targs)
			sig := e.T(fn.Type()).(*types.Signature)
			e.functionTwins(funcs, "export function "+in.name, sig, d.Type, d.Body, nil, true)
			funcs.line("")
		case *ast.TypeSpec:
			tn := in.obj.(*types.TypeName)
			named := tn.Type().(*types.Named)
			e.subst = bind(named.TypeParams(), in.targs)
			t, err := types.Instantiate(e.ctxt, named, in.targs, false)
			if err != nil {
				e.errs = append(e.errs, err)
				break
			}
			e.typeInstance(td, in.name, tn, t.(*types.Named), info)
		}
		e.info, e.subst, e.fn, e.sc = saved, savedSubst, savedFn, savedSc
	}
}

func bind(params *types.TypeParamList, args []types.Type) map[*types.TypeParam]types.Type {
	m := map[*types.TypeParam]types.Type{}
	for i := 0; i < params.Len() && i < len(args); i++ {
		m[params.At(i)] = args[i]
	}
	return m
}

// typeInstance writes the class or binding of a generic type's instance,
// with its methods, each checked against the receiver's own type
// parameters.
func (e *emitter) typeInstance(td *typeDecls, name string, tn *types.TypeName, inst *types.Named, info *types.Info) {
	var methods []*ast.FuncDecl
	for _, f := range e.filesOf(tn.Pkg()) {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv == nil {
				continue
			}
			fn, ok := info.Defs[fd.Name].(*types.Func)
			if !ok {
				continue
			}
			if namedOf(derefType(fn.Type().(*types.Signature).Recv().Type())).Obj() == tn {
				methods = append(methods, fd)
			}
		}
	}
	outer := e.subst
	e.methodSubst = func(fn *types.Func) {
		e.subst = bind(fn.Type().(*types.Signature).RecvTypeParams(), typeArgs(inst))
	}
	defer func() { e.methodSubst = nil; e.subst = outer }()
	switch u := inst.Underlying().(type) {
	case *types.Struct:
		e.structTypeDecl(td, name, inst, u, methods)
	case *types.Interface:
		td.classes.line("export type %s = any;", name)
	default:
		e.namedTypeDecl(td, name, inst, methods)
	}
}

func typeArgs(n *types.Named) []types.Type {
	var out []types.Type
	for i := 0; i < n.TypeArgs().Len(); i++ {
		out = append(out, n.TypeArgs().At(i))
	}
	return out
}

// filesOf is the syntax of a package: this module's own files when it is
// the package being emitted.
func (e *emitter) filesOf(p *types.Package) []*ast.File {
	if p == e.tpkg {
		return e.files
	}
	if q := e.prog.Packages[p.Path()]; q != nil {
		return q.Files
	}
	return nil
}
