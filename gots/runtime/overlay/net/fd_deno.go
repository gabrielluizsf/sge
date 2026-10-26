//go:build gots_overlay

// Overlay for fd_fake.go's netFD: beside the in-memory sockets Go keeps
// for js, which UDP and Unix sockets still use, a TCP connection or
// listener is one of Deno's, reached through the gotsNet bridge.
//
// Closing and setting deadlines on Deno's sockets never waits, and code
// that cannot wait does both (sync.Once closing a listener): the
// in-memory sockets, whose own may wait, are told so in a goroutine.

package net

import (
	"internal/poll"
	"syscall"
	"time"
)

type netFD struct {
	pfd poll.FD

	// immutable until Close
	family      int
	sotype      int
	isConnected bool // handshake completed or use of association with peer
	net         string
	laddr       Addr
	raddr       Addr

	// Deno's TCP connection or listener, when the netFD is one.
	sock *denoSock
	lis  *denoListener

	*fakeNetFD
}

// newDenoConn is the netFD of a Deno TCP connection.
func newDenoConn(net string, s *denoSock) *netFD {
	return &netFD{
		family:      syscall.AF_INET,
		sotype:      syscall.SOCK_STREAM,
		isConnected: true,
		net:         net,
		laddr:       tcpAddr(s.v.Call("local")),
		raddr:       tcpAddr(s.v.Call("remote")),
		sock:        s,
	}
}

func (fd *netFD) accept() (netfd *netFD, err error) {
	if fd.lis != nil {
		s, err := fd.lis.accept()
		if err != nil {
			return nil, err
		}
		return newDenoConn(fd.net, s), nil
	}
	if fd.fakeNetFD != nil {
		return fd.fakeNetFD.accept(fd.laddr)
	}
	return nil, syscall.ENOTSUP
}

func (fd *netFD) Close() error {
	switch {
	case fd.sock != nil:
		return fd.sock.close()
	case fd.lis != nil:
		return fd.lis.close()
	case fd.fakeNetFD != nil:
		go fd.fakeNetFD.Close()
		return nil
	}
	return fd.pfd.Close()
}

// shutdown has nothing to do: CloseWrite and CloseRead reach closeWrite
// and closeRead.
func (fd *netFD) shutdown(how int) error { return nil }

// closeRead leaves the socket open: Deno cannot shut reading alone, and
// what arrives is not read.
func (fd *netFD) closeRead() error { return nil }

func (fd *netFD) closeWrite() error {
	if fd.sock != nil {
		return fd.sock.closeWrite()
	}
	if fd.fakeNetFD != nil {
		go fd.fakeNetFD.closeWrite()
	}
	return nil
}

func (fd *netFD) Read(p []byte) (n int, err error) {
	if fd.sock != nil {
		return fd.sock.read(p)
	}
	if fd.fakeNetFD != nil {
		return fd.fakeNetFD.Read(p)
	}
	return 0, syscall.ENOTSUP
}

func (fd *netFD) Write(p []byte) (nn int, err error) {
	if fd.sock != nil {
		return fd.sock.write(p)
	}
	if fd.fakeNetFD != nil {
		return fd.fakeNetFD.Write(p)
	}
	return 0, syscall.ENOTSUP
}

func (fd *netFD) SetDeadline(t time.Time) error {
	switch {
	case fd.sock != nil:
		fd.sock.setReadDeadline(t)
		fd.sock.setWriteDeadline(t)
		return nil
	case fd.lis != nil:
		fd.lis.setDeadline(t)
		return nil
	case fd.fakeNetFD != nil:
		go fd.fakeNetFD.SetDeadline(t)
		return nil
	}
	return syscall.ENOTSUP
}

func (fd *netFD) SetReadDeadline(t time.Time) error {
	switch {
	case fd.sock != nil:
		fd.sock.setReadDeadline(t)
		return nil
	case fd.lis != nil:
		fd.lis.setDeadline(t)
		return nil
	case fd.fakeNetFD != nil:
		go fd.fakeNetFD.SetReadDeadline(t)
		return nil
	}
	return syscall.ENOTSUP
}

func (fd *netFD) SetWriteDeadline(t time.Time) error {
	switch {
	case fd.sock != nil:
		fd.sock.setWriteDeadline(t)
		return nil
	case fd.fakeNetFD != nil:
		go fd.fakeNetFD.SetWriteDeadline(t)
		return nil
	}
	return syscall.ENOTSUP
}
