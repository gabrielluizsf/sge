package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Package is one Go package, parsed and type-checked.
type Package struct {
	Path  string
	Name  string
	Dir   string
	Std   bool
	Files []*ast.File
	Types *types.Package
	Info  *types.Info

	// TestFiles and XTestFiles are the package's _test.go files, in the
	// package itself and in its _test package; they are only type-checked
	// when tests are transpiled.
	TestFiles  []string
	XTestFiles []string
	GoFiles    []string
	Imports    []string
	fset       *token.FileSet
}
