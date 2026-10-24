package gots

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

// emitter writes the TypeScript module of one Go package, or of a package
// together with its test files.
type emitter struct {
	prog *Program
	pkg  *Package
	info *types.Info
	tpkg *types.Package
	lay  layout
	dir  string

	imports     map[string]string // import path → module alias
	importOrder []string
	pkgNames    map[string]bool            // names taken at module level
	names       map[types.Object]string    // bindings handed out so far
	boxed       map[types.Object]bool      // variables kept in a Cell
	anon        map[string]string          // anonymous struct type → class
	anonDecls   writer                     // classes for anonymous structs
	typeConsts  map[string]string          // descriptor expression → const
	typeOrder   []string                   // typeConsts in creation order
	localTypes  map[*types.TypeName]string // types declared in functions
	fn          *funcState
	boxedDone   map[*types.Info]bool // packages findBoxed has been run on
	// predeclared are variables a goto state machine declares ahead of its
	// loop, whose Go declarations become assignments.
	predeclared   map[types.Object]bool
	publicStructs map[*types.Struct]bool   // structs with TypeScript-facing text fields
	files         []*ast.File              // the module's own files
	info0         *types.Info              // their type information
	localDecls    map[*types.Func]funcDecl // built by funcDeclOf on first use
	blk           *blocking                // built by blockingOf on first use
	bindingCall   bool                     // writing a call a defer or go statement makes later
	subst         map[*types.TypeParam]types.Type
	ctxt          *types.Context
	insts         map[string]*instance
	pendingInsts  []*instance
	methodSubst   func(fn *types.Func) // binds a generic type's method parameters
	sc            *scope
	tmp           int
	deferring     bool // writing a deferred call's arguments, which must be copies
	errs          []error
	Skipped       []error // reasons test functions were left out
	skippedTests  map[string]bool
}

// funcState is what the emitter tracks inside one function body.
type funcState struct {
	sig     *types.Signature
	ft      *ast.FuncType  // the function's own parameters and results
	async   bool           // it may block, and is an async function
	scratch string         // its temporary for calls that may hand back a promise
	body    *ast.BlockStmt // and body, which stableArg looks through
	results []string       // binding of each named result, if they are named
	frame   string         // the defer Frame, if the function defers
	breaks  []*breakable
	gotos   map[string]gotoCase // labels a goto state machine jumps to
	parent  *funcState
}

// breakable is a statement a Go break can leave: a loop, a switch or a
// select. Switches become labelled blocks, so a bare break meant for one
// has to name its label.
type breakable struct {
	label    string
	isSwitch bool
	used     bool
	// barrier is the loop a goto state machine runs in: a break or
	// continue meant for a statement outside it has to name its label.
	barrier bool
	// yield is set on a range over a function: its body is the yield
	// callback, which leaves the loop by returning.
	yield *yieldState
}

func newEmitter(prog *Program, pkg *Package, info *types.Info, tpkg *types.Package, lay layout) *emitter {
	return &emitter{
		prog: prog, pkg: pkg, info: info, tpkg: tpkg, lay: lay,
		dir:        lay.dir(pkg),
		imports:    map[string]string{},
		pkgNames:   map[string]bool{"$": true},
		names:      map[types.Object]string{},
		boxed:      map[types.Object]bool{},
		anon:       map[string]string{},
		typeConsts: map[string]string{},
		localTypes: map[*types.TypeName]string{},
		ctxt:       types.NewContext(),
		insts:      map[string]*instance{},
		info0:      info,
	}
}

func (e *emitter) errorf(pos token.Pos, format string, args ...any) {
	e.errs = append(e.errs, fmt.Errorf("%s: %s", e.prog.Fset.Position(pos), fmt.Sprintf(format, args...)))
}

// unsupported records a construct the transpiler cannot translate yet and
// returns an expression that fails loudly if it is ever reached.
func (e *emitter) unsupported(n ast.Node, what string) string {
	e.errorf(n.Pos(), "unsupported: %s", what)
	return fmt.Sprintf("$.panic(%q)", "gots: unsupported "+what)
}

func (e *emitter) newTmp() string {
	e.tmp++
	return "$" + strconv.Itoa(e.tmp)
}

func (e *emitter) typeOf(x ast.Expr) types.Type {
	if tv, ok := e.info.Types[x]; ok {
		return e.T(tv.Type)
	}
	if id, ok := x.(*ast.Ident); ok {
		if obj := e.info.ObjectOf(id); obj != nil {
			return e.T(obj.Type())
		}
	}
	return nil
}

// qualifier writes package-qualified type names the way %T prints them.
func qualifier(p *types.Package) string { return p.Name() }

func typeString(t types.Type) string { return rtString(t) }
