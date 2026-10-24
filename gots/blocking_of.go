package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// blockingAll is the analysis of every package gots writes, made once.
func (p *Program) blockingAll() *blocking {
	if p.blk == nil {
		p.blk = newBlocking(p.Fset)
		for _, path := range p.order {
			pkg := p.Packages[path]
			if pkg == nil || pkg.Info == nil || pkg.Std && !pkg.Transpile {
				continue
			}
			p.blk.addFiles(pkg.Files, pkg.Info)
		}
	}
	return p.blk
}

// blockingOf is the analysis this module is written with: the program's,
// with the module's own files added when they hold tests, checked anew.
func (e *emitter) blockingOf() *blocking {
	if e.blk != nil {
		return e.blk
	}
	e.blk = e.prog.blockingAll()
	for _, f := range e.files {
		if strings.HasSuffix(e.prog.Fset.Position(f.Pos()).Filename, "_test.go") {
			e.blk = e.blk.clone()
			e.blk.addFiles(e.files, e.info0)
			break
		}
	}
	return e.blk
}

// asyncFunc reports whether the function or closure of type ft may block.
func (e *emitter) asyncFunc(ft *ast.FuncType) bool {
	if ft == nil || !ft.Pos().IsValid() {
		return false
	}
	b := e.blockingOf()
	return b.blocks[b.key(ft.Pos())]
}

// blocks reports whether calling fn may block.
func (e *emitter) blocks(fn *types.Func) bool {
	if blockingNatives[nativeName(fn)] {
		return true
	}
	s := fn.Origin().Scope()
	if s == nil || !s.Pos().IsValid() {
		return false
	}
	b := e.blockingOf()
	return b.blocks[b.key(s.Pos())]
}
