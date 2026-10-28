//go:build gots_overlay

// Overlay for the assembly ChaCha8 block function: the generic Go one.

package chacha8rand

func block(seed *[4]uint64, blocks *[32]uint64, counter uint32) {
	block_generic(seed, blocks, counter)
}
