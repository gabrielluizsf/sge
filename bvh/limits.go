package bvh

// Limits bound what a file may make a reader do: how many bytes it may
// have, joints it may name, how deep they may nest, and frames it may
// hold.
type Limits struct {
	Bytes                 int64
	Joints, Depth, Frames int
}

// DefaultLimits are ample for any skeleton and a few minutes of motion.
var DefaultLimits = Limits{Bytes: 16 << 20, Joints: 1024, Depth: 64, Frames: 1 << 14}
