//go:build gots_overlay

// Overlay for Unix datagram sockets: Deno's, through the gotsNet bridge,
// behind --unstable-net as UDP is. Deno binds every such socket to a
// path; one dialed without a local address, which the system would leave
// unnamed, is bound to a fresh one in the temporary directory, removed
// when it closes.

package net

import (
	"context"
	"os"
	"strconv"
	"syscall"
	"syscall/js"
	"time"
)

func (sl *sysListener) listenUnixgram(ctx context.Context, laddr *UnixAddr) (*UnixConn, error) {
	if c := sl.ListenConfig.Control; c != nil {
		if err := c("unixgram", laddr.Name, denoRawConn{}); err != nil {
			return nil, err
		}
	}
	return openUnixgram(laddr, nil, false)
}

func (sd *sysDialer) dialUnixgram(ctx context.Context, laddr, raddr *UnixAddr) (*UnixConn, error) {
	ctrl := sd.Dialer.ControlContext
	if ctrl == nil && sd.Dialer.Control != nil {
		ctrl = func(ctx context.Context, network, address string, c syscall.RawConn) error {
			return sd.Dialer.Control(network, address, c)
		}
	}
	if ctrl != nil {
		if err := ctrl(ctx, "unixgram", raddr.Name, denoRawConn{}); err != nil {
			return nil, err
		}
	}
	unnamed := laddr == nil
	if unnamed {
		laddr = &UnixAddr{Name: tempSocketPath(), Net: "unixgram"}
	}
	return openUnixgram(laddr, raddr, unnamed)
}

// openUnixgram opens Deno's Unix datagram socket at laddr, connected to
// raddr if given; unlink has Close remove its file.
func openUnixgram(laddr, raddr *UnixAddr, unlink bool) (*UnixConn, error) {
	v := js.Global().Get("gotsNet").Call("listenUnixgram", laddr.Name)
	if v.Get("failed").Truthy() {
		if v.Get("name").String() == "NotSupported" {
			return nil, os.NewSyscallError("socket", syscall.ENOTSUP)
		}
		return nil, denoError(v)
	}
	d := &denoDgram{v: v}
	if unlink {
		d.unlink = laddr.Name
	}
	fd := &netFD{family: syscall.AF_UNIX, sotype: syscall.SOCK_DGRAM, net: "unixgram", laddr: laddr, isConnected: raddr != nil, dgram: d}
	if raddr != nil {
		d.peer = &dgramAddr{raddr.Name, -1}
		fd.raddr = raddr
	}
	return newUnixConn(fd), nil
}

// tempSocketPath is a path no socket has, for an unnamed one.
func tempSocketPath() string {
	return os.TempDir() + "/gots-unixgram-" + strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.FormatUint(runtime_rand()%1e9, 36)
}

// readFrom and writeTo are a Unix datagram socket's; any other netFD's
// are the in-memory network's.

func (fd *netFD) readFrom(p []byte) (int, syscall.Sockaddr, error) {
	if fd.dgram == nil || fd.family != syscall.AF_UNIX {
		return fd.fakeNetFD.readFrom(p)
	}
	n, from, err := fd.dgram.receive(p)
	if err != nil {
		return n, nil, err
	}
	return n, &syscall.SockaddrUnix{Name: from.host}, nil
}

func (fd *netFD) writeTo(p []byte, sa syscall.Sockaddr) (int, error) {
	if fd.dgram == nil || fd.family != syscall.AF_UNIX {
		return fd.fakeNetFD.writeTo(p, sa)
	}
	u, ok := sa.(*syscall.SockaddrUnix)
	if !ok {
		return 0, syscall.EAFNOSUPPORT
	}
	return fd.dgram.send(p, &dgramAddr{u.Name, -1})
}
