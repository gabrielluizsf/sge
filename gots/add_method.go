package gots

import (
	"go/ast"
	"go/types"
)

// addMethod notes that the method fd is one a call through an interface
// may reach: if it may block, so may every such call.
func (b *blocking) addMethod(fd *ast.FuncDecl, info *types.Info) {
	fn, ok := info.Defs[fd.Name].(*types.Func)
	if !ok {
		return
	}
	me, k := b.key(fd.Type.Pos()), ifaceKey(fn)
	b.callers[me] = append(b.callers[me], k)
	if b.blocks[me] {
		b.mark(k)
	}
}
