package gots

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"
)

// typeDecls are what a package's type declarations turn into, kept apart so
// the module can put classes before the descriptors that refer to them.
type typeDecls struct {
	classes writer // classes, interfaces and type aliases
	descs   writer // descriptors of structs and interfaces
	funcs   writer // methods of types that are not structs
	// named are the descriptors of other named types, which are built
	// from their underlying types' at once and so go in dependency order.
	named []namedDesc
}

type namedDesc struct {
	name  string   // the type's binding
	lines []string // its descriptor and registrations
	deps  []string // bindings of the package's named types it is made of
}

// writeNamed writes the named type descriptors, each after those it uses.
func (td *typeDecls) writeNamed(w *writer) {
	byName := map[string]*namedDesc{}
	for i := range td.named {
		byName[td.named[i].name] = &td.named[i]
	}
	done := map[string]bool{}
	var visit func(d *namedDesc)
	visit = func(d *namedDesc) {
		if done[d.name] {
			return
		}
		done[d.name] = true
		for _, dep := range d.deps {
			if o := byName[dep]; o != nil {
				visit(o)
			}
		}
		for _, l := range d.lines {
			w.line("%s", l)
		}
	}
	for i := range td.named {
		visit(&td.named[i])
	}
}

// methodsOf groups a package's method declarations by receiver type.
func (e *emitter) methodsOf(files []*ast.File) map[*types.TypeName][]*ast.FuncDecl {
	out := map[*types.TypeName][]*ast.FuncDecl{}
	for _, f := range files {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv == nil {
				continue
			}
			fn := e.info.Defs[fd.Name].(*types.Func)
			n := namedOf(derefType(e.recvType(fn)))
			out[n.Obj()] = append(out[n.Obj()], fd)
		}
	}
	return out
}

func (e *emitter) typeDecl(td *typeDecls, tn *types.TypeName, methods []*ast.FuncDecl) {
	name := mangle(tn.Name())
	if tn.IsAlias() {
		target := types.Unalias(tn.Type())
		td.classes.line("export type %s = %s;", name, e.tsType(target))
		if isStruct(target) {
			if _, named := target.(*types.Named); named {
				// With the descriptors, after every class: an instance of a
				// generic struct has its class written after the declarations.
				td.descs.line("export const %s = %s;", name, e.structClass(target))
			}
		}
		return
	}
	t := tn.Type().(*types.Named)
	switch u := t.Underlying().(type) {
	case *types.Struct:
		e.structTypeDecl(td, name, t, u, methods)
	case *types.Interface:
		e.ifaceTypeDecl(td, name, t, u)
	default:
		e.namedTypeDecl(td, name, t, methods)
	}
}

// structTypeDecl writes the class of a named struct with its methods, the
// methods promoted from embedded fields, and the class an interface keeps
// its values in.
func (e *emitter) structTypeDecl(td *typeDecls, name string, t *types.Named, st *types.Struct, methods []*ast.FuncDecl) {
	w := &td.classes
	e.structClassDecl(w, name, st, true, func() {
		for _, fd := range methods {
			e.method(w, &td.funcs, fd, true)
		}
		e.promoted(w, types.NewMethodSet(types.NewPointer(t)), "this")
	})
	valCls := "null"
	vms := types.NewMethodSet(t)
	if vms.Len() > 0 {
		valCls = name + "$Val"
		w.open("export class %s extends $.Val {", valCls)
		for i := 0; i < vms.Len(); i++ {
			fn := vms.At(i).Obj().(*types.Func)
			m := fn.Name()
			w.line("%s(...a: any[]): any { return this.v.%s(...a); }", fieldName(m), fieldName(m))
			switch {
			case e.twin(fn):
				w.line("%s$a(...a: any[]): any { return this.v.%s$a(...a); }", fieldName(m), e.goName(fn))
			case ifaceMethod(fn):
				w.line("%s$a(...a: any[]): any { return $.callA(this.v, %q, ...a); }", fieldName(m), fieldName(m))
			}
		}
		w.close("}")
	}
	desc := fmt.Sprintf("$.structType(%q, true, %s, %s, () => %s)", typeString(t), name, valCls, e.fieldDescs(st))
	td.descs.line("%s.$type = %s;", name, e.canonical(t, e.inPkg(t, desc)))
	if ms := e.methodDescs(t); ms != "" {
		td.descs.line("%s.$type.setMethods(() => %s);", name, ms)
	}
}

// promoted writes forwarding methods for methods a struct gets from its
// embedded fields.
func (e *emitter) promoted(w *writer, ms *types.MethodSet, self string) {
	for i := 0; i < ms.Len(); i++ {
		sel := ms.At(i)
		if len(sel.Index()) == 1 {
			continue
		}
		fn := sel.Obj().(*types.Func)
		path, t := e.fieldPath(self, e.T(sel.Recv()), sel.Index()[:len(sel.Index())-1])
		m := fieldName(fn.Name())
		if onClass(fn) {
			w.line("%s(...a: any[]): any { return %s.%s(...a); }", m, path, m)
			if g := e.goName(fn); g != m {
				w.line("%s(...a: any[]): any { return %s.%s(...a); }", g, path, g)
			}
			switch {
			case e.twin(fn):
				w.line("%s$a(...a: any[]): any { return %s.%s$a(...a); }", e.goName(fn), path, e.goName(fn))
			case ifaceMethod(fn):
				// What the field holds may have a twin: code that can wait
				// reaches it through this one.
				w.line("%s$a(...a: any[]): any { return $.callA(%s, %q, ...a); }", m, path, m)
			}
			continue
		}
		recv := path
		_, ptrRecv := e.recvType(fn).(*types.Pointer)
		switch {
		case ptrRecv && !isPointer(t):
			dot := strings.LastIndex(path, ".")
			recv = "new $.FieldPtr(" + path[:dot] + ", " + jsQuote(path[dot+1:]) + ")"
		case !ptrRecv && isPointer(t):
			recv = path + ".v"
		}
		w.line("%s(...a: any[]): any { return %s(%s, ...a); }", m, e.methodFunc(fn), recv)
		if e.twin(fn) {
			w.line("%s$a(...a: any[]): any { return %s$a(%s, ...a); }", m, e.methodFunc(fn), recv)
		}
	}
}

func (e *emitter) ifaceTypeDecl(td *typeDecls, name string, t *types.Named, u *types.Interface) {
	w := &td.classes
	if u.Empty() {
		w.line("export type %s = any;", name)
	} else if !u.IsMethodSet() {
		w.line("export type %s = any;", name)
	} else {
		w.open("export interface %s {", name)
		for i := 0; i < u.NumMethods(); i++ {
			m := u.Method(i)
			w.line("%s%s;", fieldName(m.Name()), e.tsMethodType(m.Type().(*types.Signature)))
		}
		// A Go interface asks for these methods and no fewer, never for no
		// more: an object with others (as an optional interface, say Painter
		// beside Script, looks for) still is one.
		w.line("[more: string]: any;")
		w.close("}")
	}
	td.descs.line("export const %s$type: $.Type = %s;", name, e.inPkg(t, fmt.Sprintf("$.ifaceType(%q, %s, %s)", typeString(t), methodList(u), e.methodSigs(u))))
}

// namedTypeDecl writes a named type that is not a struct: a type alias for
// TypeScript, its methods as functions taking the receiver first, and the
// classes interfaces keep its values and pointers to it in.
func (e *emitter) namedTypeDecl(td *typeDecls, name string, t *types.Named, methods []*ast.FuncDecl) {
	td.classes.line("export type %s = %s;", name, e.tsType(t.Underlying()))
	for _, fd := range methods {
		e.method(&td.funcs, &td.funcs, fd, false)
	}
	valCls := "null"
	if vms := types.NewMethodSet(t); vms.Len() > 0 {
		valCls = name + "$Val"
		e.valClass(&td.classes, valCls, vms, "this.v")
	}
	desc := fmt.Sprintf("$.named(%q, %s, %s)", typeString(t), e.descInline(t.Underlying()), valCls)
	nd := namedDesc{name: name, deps: e.namedDeps(t.Underlying())}
	nd.lines = append(nd.lines, fmt.Sprintf("export const %s$type: $.Type = %s;", name, e.canonical(t, e.inPkg(t, desc))))
	if ms := e.methodDescs(t); ms != "" {
		nd.lines = append(nd.lines, fmt.Sprintf("%s$type.setMethods(() => %s);", name, ms))
	}
	if pms := types.NewMethodSet(types.NewPointer(t)); pms.Len() > 0 {
		e.valClass(&td.classes, name+"$PtrVal", pms, "this.v")
		nd.lines = append(nd.lines, fmt.Sprintf("$.ptrTo(%s$type).valCls = %s$PtrVal;", name, name))
	}
	td.named = append(td.named, nd)
}

// valClass writes the Val subclass carrying the method set ms for values
// that are not objects.
func (e *emitter) valClass(w *writer, cls string, ms *types.MethodSet, v string) {
	w.open("export class %s extends $.Val {", cls)
	for i := 0; i < ms.Len(); i++ {
		fn := ms.At(i).Obj().(*types.Func)
		_, ptrRecv := e.recvType(fn).(*types.Pointer)
		recv := v
		if !ptrRecv && ms.At(i).Indirect() {
			recv = v + ".v"
		}
		_, isPtrSet := ms.At(i).Recv().(*types.Pointer)
		if isPtrSet && !ptrRecv {
			recv = v + ".v"
		}
		w.line("%s(...a: any[]): any { return %s(%s, ...a); }", fieldName(fn.Name()), e.methodFunc(fn), recv)
		if e.twin(fn) {
			w.line("%s$a(...a: any[]): any { return %s$a(%s, ...a); }", fieldName(fn.Name()), e.methodFunc(fn), recv)
		}
	}
	w.close("}")
}

// sortedTypes are the package's type names, in declaration order.
func sortedTypes(scope *types.Scope) []*types.TypeName {
	var out []*types.TypeName
	for _, n := range scope.Names() {
		if tn, ok := scope.Lookup(n).(*types.TypeName); ok {
			out = append(out, tn)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Pos() < out[j].Pos() })
	return out
}

// canonical wraps the descriptor of a generic type's instance so every
// module instantiating it shares one: instances are emitted by each module
// that uses them, but they are one type.
func (e *emitter) canonical(t *types.Named, desc string) string {
	if t.TypeArgs().Len() == 0 {
		return desc
	}
	return fmt.Sprintf("$.canonical(%q, () => %s)", types.TypeString(t, nil), desc)
}

// namedDeps are the bindings of this module's named non-struct types that
// t is built from directly, not through a struct (whose fields are read
// lazily).
func (e *emitter) namedDeps(t types.Type) []string {
	var out []string
	var walk func(t types.Type)
	walk = func(t types.Type) {
		switch t := types.Unalias(t).(type) {
		case *types.Named:
			if t.Obj().Pkg() == e.tpkg && !isStruct(t) && !isInterface(t) {
				out = append(out, e.namedRef(t))
			}
		case *types.Pointer:
			walk(t.Elem())
		case *types.Slice:
			walk(t.Elem())
		case *types.Array:
			walk(t.Elem())
		case *types.Map:
			walk(t.Key())
			walk(t.Elem())
		case *types.Chan:
			walk(t.Elem())
		}
	}
	walk(t)
	return out
}

// ifaceMethod reports whether fn is a method of an interface, as one
// promoted from an embedded interface field is.
func ifaceMethod(fn *types.Func) bool {
	r := fn.Type().(*types.Signature).Recv()
	return r != nil && types.IsInterface(r.Type())
}
