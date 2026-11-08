//go:build linux && amd64

package gl

import (
	"github.com/gabrielluizsf/sge/native"
)

// status is how a shader's or a program's state is asked: get, its
// glGet*iv, log, its glGet*InfoLog, of the object id, for what.
type status struct {
	get, log native.Proc
	id       uint32
	what     uint32
}

// status is whether s's object succeeded at what s asks.
func (d *Device) status(s status) bool {
	var ok int32
	call(s.get, n(s.id), n(s.what), ptr(pointer(&ok)))
	return ok != 0
}

// infoLog is what GL says of s's object.
func (d *Device) infoLog(s status) string {
	var size int32
	call(s.get, n(s.id), n(glInfoLogLength), ptr(pointer(&size)))
	if size <= 0 {
		return "(no log)"
	}
	buf := make([]byte, size)
	call(s.log, n(s.id), n(size), n(0), ptr(pointer(&buf[0])))
	return string(buf[:len(buf)-1])
}
