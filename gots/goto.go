package gots

import (
	"go/ast"
	"go/token"
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
