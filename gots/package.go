package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Package is one Go package, parsed and type-checked.
type Package struct {
	Path string
	Name string
	Dir  string
	Std  bool
	// Transpile is set for a standard library package that has no
	// hand-written TypeScript and is transpiled from its source.
	Transpile bool
	Files     []*ast.File
	Types     *types.Package
	Info      *types.Info

	// TestFiles and XTestFiles are the package's _test.go files, in the
	// package itself and in its _test package; they are only type-checked
	// when tests are transpiled.
	TestFiles  []string
	XTestFiles []string
	GoFiles    []string
	Imports    []string
	// ImportMap maps an import path as written to the package it means,
	// for the standard library's vendored packages.
	ImportMap map[string]string
	fset      *token.FileSet
}
