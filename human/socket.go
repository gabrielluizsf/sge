package human

import "github.com/gabrielluizsf/sge/gltf"

// Socket is a place on the body where a thing is carried or worn, following
// the bone it is on: a hat on the head, a sword at the hip, a pack on the
// back. Offset is from the bone, in metres, in the bone's frame, turned by
// Turn, a quaternion.
type Socket struct {
	Bone   string
	Offset [3]float64
	Turn   [4]float64
}

// The sockets a game asks for most.
var (
	SocketHead      = Socket{Bone: "head", Offset: [3]float64{0, 0.14, 0}, Turn: [4]float64{3: 1}}
	SocketBack      = Socket{Bone: "spine02", Offset: [3]float64{0, 0.05, -0.16}, Turn: [4]float64{3: 1}}
	SocketHipLeft   = Socket{Bone: "pelvis.L", Offset: [3]float64{0.1, 0, 0}, Turn: [4]float64{3: 1}}
	SocketHipRight  = Socket{Bone: "pelvis.R", Offset: [3]float64{-0.1, 0, 0}, Turn: [4]float64{3: 1}}
	SocketHandLeft  = Socket{Bone: "wrist.L", Offset: [3]float64{0, -0.09, 0.01}, Turn: [4]float64{3: 1}}
	SocketHandRight = Socket{Bone: "wrist.R", Offset: [3]float64{0, -0.09, 0.01}, Turn: [4]float64{3: 1}}
)

// Sockets are all of them.
var Sockets = []Socket{SocketHead, SocketBack, SocketHipLeft, SocketHipRight, SocketHandLeft, SocketHandRight}

// Grip is how a thing is held at the socket, through the animations
// during: on its one bone.
func (s Socket) Grip(during ...string) Grip {
	return Grip{Hands: [2]string{s.Bone, s.Bone}, During: during, Offset: s.Offset, Turn: s.Turn}
}

// Attach gives person m thing at the socket, through the animations during,
// and out of sight in the others. m's animations are made first.
func Attach(m, thing *gltf.Model, s Socket, during ...string) error {
	return Hold(m, thing, s.Grip(during...))
}
