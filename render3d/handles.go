package render3d

// Buffer, Texture and Pipeline are a device's resources, by number: 0 is
// none, and a backend numbers its own from 1.
type (
	Buffer   uint32
	Texture  uint32
	Pipeline uint32
)

// Screen is the texture a pass draws into, and Read reads, for the
// device's own screen or window.
const Screen Texture = 0

// Handle is any of a device's resources, for Release.
type Handle interface{ isHandle() }

func (Buffer) isHandle()   {}
func (Texture) isHandle()  {}
func (Pipeline) isHandle() {}
