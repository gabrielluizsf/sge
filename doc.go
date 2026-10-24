// Package sge is a small game engine written in Go on top of antui, that
// also ships to Deno: gots transpiles its packages to TypeScript, so the
// same engine is a Go module and a JSR package.
//
// The engine is in its packages: vector (3D math), transform (position,
// rotation, scale), collider (overlap tests), clock (frame time), sprite
// (antui canvases drawn with a pivot, scale and rotation, from PNG or
// SVG), animation (frames by name), game (the world, its objects and
// scripts, the camera and the window loop), tilemap (grids of tiles,
// animated, walked through) and terrain (endless tile maps of named
// areas, made around the camera).
//
// go generate writes the TypeScript to deno/.
package sge

//go:generate go run ./cmd/gots -o deno -index ./vector ./transform ./collider ./clock ./sprite ./animation ./game ./tilemap ./terrain github.com/gabrielluizsf/antui github.com/gabrielluizsf/antui/canvas
