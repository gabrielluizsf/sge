// Package bvh reads Biovision Hierarchy motion: a skeleton of joints,
// each offset from its parent and moved by its channels, and frames of
// those channels' values, as MakeHuman's pose units and motion capture
// are kept.
//
// A file is read as untrusted, after CVE-2025-6119, a use after free in
// Assimp's BVH loader on a crafted file: Go keeps memory safe, and what is
// left, a file that would take all memory or time, is bounded by Limits:
// its size, its joints and how deep they nest, read without recursion,
// and its frames, whose values must be as many as their channels and
// finite.
package bvh
