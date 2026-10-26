//go:build gots_overlay

// The netFD methods UDPConn reads and writes through, for one of Deno's
// UDP sockets; for any other netFD, the in-memory network's.

package net

import "syscall"

func (fd *netFD) readFromInet4(p []byte, from *syscall.SockaddrInet4) (int, error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.readFromInet4(p, from)
	}
	n, a, err := fd.dgram.receive(p)
	if err == nil {
		from.Port = a.Port
		copy(from.Addr[:], a.IP.To4())
	}
	return n, err
}

func (fd *netFD) readFromInet6(p []byte, from *syscall.SockaddrInet6) (int, error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.readFromInet6(p, from)
	}
	n, a, err := fd.dgram.receive(p)
	if err == nil {
		from.Port = a.Port
		copy(from.Addr[:], a.IP.To16())
	}
	return n, err
}

func (fd *netFD) readMsgInet4(p []byte, oob []byte, flags int, sa *syscall.SockaddrInet4) (n, oobn, retflags int, err error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.readMsgInet4(p, oob, flags, sa)
	}
	n, err = fd.readFromInet4(p, sa)
	return n, 0, 0, err
}

func (fd *netFD) readMsgInet6(p []byte, oob []byte, flags int, sa *syscall.SockaddrInet6) (n, oobn, retflags int, err error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.readMsgInet6(p, oob, flags, sa)
	}
	n, err = fd.readFromInet6(p, sa)
	return n, 0, 0, err
}

func (fd *netFD) writeToInet4(p []byte, sa *syscall.SockaddrInet4) (int, error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.writeToInet4(p, sa)
	}
	return fd.dgram.send(p, &UDPAddr{IP: IP(sa.Addr[:]), Port: sa.Port})
}

func (fd *netFD) writeToInet6(p []byte, sa *syscall.SockaddrInet6) (int, error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.writeToInet6(p, sa)
	}
	return fd.dgram.send(p, &UDPAddr{IP: IP(sa.Addr[:]), Port: sa.Port})
}

func (fd *netFD) writeMsgInet4(p []byte, oob []byte, sa *syscall.SockaddrInet4) (n int, oobn int, err error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.writeMsgInet4(p, oob, sa)
	}
	n, err = fd.writeToInet4(p, sa)
	return n, 0, err
}

func (fd *netFD) writeMsgInet6(p []byte, oob []byte, sa *syscall.SockaddrInet6) (n int, oobn int, err error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.writeMsgInet6(p, oob, sa)
	}
	n, err = fd.writeToInet6(p, sa)
	return n, 0, err
}

func (fd *netFD) writeMsg(p []byte, oob []byte, sa syscall.Sockaddr) (n int, oobn int, err error) {
	if fd.dgram == nil {
		return fd.fakeNetFD.writeMsg(p, oob, sa)
	}
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return fd.writeMsgInet4(p, oob, sa)
	case *syscall.SockaddrInet6:
		return fd.writeMsgInet6(p, oob, sa)
	}
	return 0, 0, syscall.EAFNOSUPPORT
}
