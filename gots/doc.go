// Package gots transpiles Go packages to TypeScript that runs on Deno.
//
// The Go type checker does the hard work: every package is loaded and checked
// for GOOS=js GOARCH=wasm, so build tags pick the pure-Go code paths, and the
// emitter only translates what the checker already proved correct. The
// TypeScript it writes leans on a small runtime (runtime/rt) for what
// JavaScript lacks: slices, pointers to variables, Go's integer widths,
// interfaces holding non-objects, and panics. Standard library packages are
// hand-written TypeScript with Go's signatures (runtime/std).
//
// How Go values map to JavaScript:
//
//   - Structs are class instances. A pointer to a struct is the instance
//     itself, so copying a struct value is an explicit $clone and assigning
//     one is $set, in place, keeping pointers to it valid.
//   - Arrays are JavaScript or typed arrays, copied like structs.
//   - Slices are rt.Slice windows over such arrays; nil is null.
//   - Pointers to anything else are rt.Ptr objects with a v property: a Cell
//     for a variable whose address is taken, accessors for fields and
//     elements.
//   - int, int64, uint and float64 are numbers; the narrower integers wrap on
//     every operation; uint64 is a bigint.
//   - Strings are byte strings: one UTF-16 code unit per byte.
//   - Interfaces hold struct pointers as the object and anything else in an
//     rt.Val carrying the dynamic type and its methods.
package gots
