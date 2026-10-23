package gots

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"
)

// gotoRegion is a run of statements [start, end) that a forward goto
// leaves: the statements become a labelled block, the goto a break out of
// it, which lands on the label at end.
type gotoRegion struct {
	start, end int
	label      string
}

// gotoTargets lists the labels in list that a statement before them jumps
// to with goto. Go forbids jumping over declarations still in scope at
// the label, so wrapping the statements in between in a block is safe.
func gotoRegions(list []ast.Stmt) []gotoRegion {
	var regions []gotoRegion
	for i, s := range list {
		ls, ok := s.(*ast.LabeledStmt)
		if !ok {
			continue
		}
		start := -1
		for j := 0; j < i; j++ {
			if jumpsTo(list[j], ls.Label.Name) {
				start = j
				break
			}
		}
		if start >= 0 {
			regions = append(regions, gotoRegion{start, i, ls.Label.Name})
		}
	}
	sort.Slice(regions, func(a, b int) bool { return regions[a].end < regions[b].end })
	// Overlapping regions are made to nest: a later one starts no later
	// than any earlier one it overlaps.
	for k := range regions {
		for m := 0; m < k; m++ {
			if regions[m].end > regions[k].start && regions[m].start < regions[k].start {
				regions[k].start = regions[m].start
			}
		}
	}
	return regions
}

// jumpsTo reports whether s contains goto label, outside closures.
func jumpsTo(s ast.Stmt, label string) bool {
	found := false
	ast.Inspect(s, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.BranchStmt:
			if n.Tok == token.GOTO && n.Label.Name == label {
				found = true
			}
		}
		return !found
	})
	return found
}

func gotoLabel(name string) string { return mangle(name) + "$g" }

// stmtsWithGotos writes list, opening a labelled block where each goto
// region starts and closing it where its label is.
func (e *emitter) stmtsWithGotos(w *writer, list []ast.Stmt, regions []gotoRegion) {
	for i, s := range list {
		for k := len(regions) - 1; k >= 0; k-- {
			if regions[k].start == i {
				w.open("%s: {", gotoLabel(regions[k].label))
			}
		}
		for _, r := range regions {
			if r.end == i {
				w.close("}")
			}
		}
		e.stmt(w, s)
	}
}

// gotoCase is where a goto lands in a state machine: the case of the
// switch on state, in the loop that runs it.
type gotoCase struct {
	state, loop string
	index       int
}

// backwardGoto reports whether a label in list is the target of a goto in
// a statement at or after it.
func backwardGoto(list []ast.Stmt) bool {
	for i, s := range list {
		ls, ok := s.(*ast.LabeledStmt)
		if !ok {
			continue
		}
		for j := i; j < len(list); j++ {
			if jumpsTo(list[j], ls.Label.Name) {
				return true
			}
		}
	}
	return false
}

// stmtsStateMachine writes a statement list whose labels gotos jump to in
// both directions: the list becomes a switch on a state variable, each
// target label opening a case, inside a loop a goto continues after
// setting the state. Cases fall through, so execution runs on from one
// label into the next as it does in Go.
func (e *emitter) stmtsStateMachine(w *writer, list []ast.Stmt) {
	state, loop := e.sc.fresh("$g"), e.sc.fresh("$gm")
	if e.fn.gotos == nil {
		e.fn.gotos = map[string]gotoCase{}
	}
	cases := map[int]int{}
	for i, s := range list {
		ls, ok := s.(*ast.LabeledStmt)
		if !ok {
			continue
		}
		for j := range list {
			if jumpsTo(list[j], ls.Label.Name) {
				cases[i] = len(cases) + 1
				e.fn.gotos[ls.Label.Name] = gotoCase{state, loop, cases[i]}
				break
			}
		}
	}
	w.line("let %s = 0;", state)
	// Each pass of the loop enters the switch's block afresh, so a jump
	// past a declaration would leave the variable uninitialized: the
	// list's own variables are declared before the loop instead.
	if e.predeclared == nil {
		e.predeclared = map[types.Object]bool{}
	}
	for _, s := range list {
		for _, v := range declaredBy(e.info, s) {
			n := e.declare(v)
			if e.boxed[v] {
				w.line("let %s: $.Cell<%s> = null;", n, e.tsType(e.varType(v)))
			} else {
				w.line("let %s: %s = %s;", n, e.tsType(e.varType(v)), e.zero(e.varType(v)))
			}
			e.predeclared[v] = true
		}
	}
	e.fn.breaks = append(e.fn.breaks, &breakable{barrier: true})
	w.open("%s: for (;;) {", loop)
	w.open("switch (%s) {", state)
	w.line("case 0:")
	for i, s := range list {
		if k, ok := cases[i]; ok {
			w.line("case %d:", k)
		}
		w.indent++
		e.stmt(w, s)
		w.indent--
	}
	w.close("}")
	w.line("break;")
	w.close("}")
	e.fn.breaks = e.fn.breaks[:len(e.fn.breaks)-1]
}

// declaredBy are the variables statement s declares in the block it is in
// (not in blocks of its own), looking through labels.
func declaredBy(info *types.Info, s ast.Stmt) []*types.Var {
	var out []*types.Var
	add := func(id *ast.Ident) {
		if v, ok := info.Defs[id].(*types.Var); ok && id.Name != "_" {
			out = append(out, v)
		}
	}
	switch s := s.(type) {
	case *ast.LabeledStmt:
		return declaredBy(info, s.Stmt)
	case *ast.AssignStmt:
		if s.Tok == token.DEFINE {
			for _, l := range s.Lhs {
				if id, ok := l.(*ast.Ident); ok {
					add(id)
				}
			}
		}
	case *ast.DeclStmt:
		if g, ok := s.Decl.(*ast.GenDecl); ok && g.Tok == token.VAR {
			for _, sp := range g.Specs {
				for _, id := range sp.(*ast.ValueSpec).Names {
					add(id)
				}
			}
		}
	}
	return out
}
