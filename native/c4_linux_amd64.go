package native

// hookFrame is a C call with at most four integer arguments. The hooks
// keep one per function and reuse it, to stay inside the linker's
// 800-byte nosplit budget.
type hookFrame struct {
	fn  uintptr
	a   [4]uintptr
	ret uintptr
}

// hookCall calls f.fn with the System V ABI on the current stack, which
// must be a C (system) stack. Used only by the runtime hooks.
//
//go:noescape
func hookCall(f *hookFrame)

//go:nosplit
//go:norace
func (f *hookFrame) c(fn uintptr, a [4]uintptr) uintptr {
	f.fn, f.a = fn, a
	hookCall(f)
	return f.ret
}
