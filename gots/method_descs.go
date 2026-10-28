package gots

import (
	"fmt"
	"go/types"
	"strings"
)

// methodDescs is the exported method set of the named type t for its
// descriptor, as reflect sees it: the methods of *t, promoted ones
// included, each with its type without the receiver and whether t itself
// has it. It is "" for a type with none.
func (e *emitter) methodDescs(t *types.Named) string {
	ptr := types.NewMethodSet(types.NewPointer(t))
	val := types.NewMethodSet(t)
	var out []string
	for i := 0; i < ptr.Len(); i++ {
		sel := ptr.At(i)
		if !sel.Obj().Exported() {
			continue
		}
		inVal := val.Lookup(sel.Obj().Pkg(), sel.Obj().Name()) != nil
		out = append(out, fmt.Sprintf("[%q, %s, %v]", sel.Obj().Name(), e.descInline(e.T(sel.Type())), inVal))
	}
	if len(out) == 0 {
		return ""
	}
	return "[" + strings.Join(out, ", ") + "]"
}
