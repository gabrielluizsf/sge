package gots

import (
	"go/ast"
	"go/types"
	"strings"
)

// blockingAll is the analysis of every package gots writes, made once: of
// what blocks for sure, or with cha of what may block through interfaces.
func (p *Program) blockingAll(cha bool) *blocking {
	at := &p.blk
	if cha {
		at = &p.blkMay
	}
	if *at == nil {
		*at = newBlocking(p.Fset, cha)
		for _, path := range p.order {
			pkg := p.Packages[path]
			if pkg == nil || pkg.Info == nil || pkg.Std && !pkg.Transpile {
				continue
			}
			(*at).addFiles(pkg.Files, pkg.Info)
		}
	}
	return *at
}

// blockingOf is the analysis this module is written with: the program's,
// with the module's own files added when they hold tests, checked anew.
func (e *emitter) blockingOf() *blocking {
	if e.blk == nil {
		e.blk = e.withTests(e.prog.blockingAll(false))
	}
	return e.blk
}

// mayBlockOf is blockingOf for what may block through interfaces.
func (e *emitter) mayBlockOf() *blocking {
	if e.blkMay == nil {
		e.blkMay = e.withTests(e.prog.blockingAll(true))
	}
	return e.blkMay
}

// withTests is b, or a copy of it with the module's files added when they
// hold tests, checked anew.
func (e *emitter) withTests(b *blocking) *blocking {
	for _, f := range e.files {
		if strings.HasSuffix(e.prog.Fset.Position(f.Pos()).Filename, "_test.go") {
			b = b.clone()
			b.addFiles(e.files, e.info0)
			break
		}
	}
	return b
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
