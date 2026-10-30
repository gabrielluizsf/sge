// Package sge is a small game engine written in Go on top of antui, that
// also ships to Deno: gots transpiles its packages to TypeScript, so the
// same engine is a Go module and a JSR package.
//
// The engine is in its packages: vector (3D math), transform (position,
// rotation, scale), collider (overlap tests), clock (frame time), sprite
// (antui canvases drawn with a pivot, scale and rotation, from PNG or
// SVG), animation (frames by name), game (the world, its objects and
// scripts, the camera and the window loop), tilemap (grids of tiles,
// animated, walked through), terrain (endless tile maps of named areas,
// made around the camera), iso (an isometric view), sound (a mixer of
// clips and loops), town (towns laid out over an endless world), pathfind
// (ways across a grid), daytime (the hours, day and night), scene
// (switching between scenes), inventory (bags of items), request (orders
// of random items), random (a sequence the same everywhere), label
// (text drawn smooth and outlined), gltf (3D models read as untrusted),
// bake (models rendered into isometric sprites), loader (heavy assets
// loaded a few milliseconds a frame) and fetch (assets downloaded and
// checked against their pins).
//
// go generate writes the TypeScript to deno/.
package sge

//go:generate go run ./cmd/gots -o deno -index ./vector ./transform ./collider ./clock ./sprite ./animation ./game ./tilemap ./terrain ./iso ./sound ./town ./pathfind ./daytime ./scene ./inventory ./request ./random ./label ./gltf ./bake ./loader ./fetch github.com/gabrielluizsf/antui github.com/gabrielluizsf/antui/canvas
