//go:build gots_overlay

// Overlay for Unix stream sockets: Deno's, through the gotsNet bridge.
// Deno has no datagram or packet Unix sockets; those stay Go's in-memory
// ones for js.

package net

import (
	"context"
	"os"
	"syscall"
	"syscall/js"
)

func (sd *sysDialer) dialUnix(ctx context.Context, laddr, raddr *UnixAddr) (*UnixConn, error) {
	if sd.network != "unix" {
		return nil, os.NewSyscallError("socket", syscall.ENOTSUP)
	}
	ctrl := sd.Dialer.ControlContext
	if ctrl == nil && sd.Dialer.Control != nil {
		ctrl = func(ctx context.Context, network, address string, c syscall.RawConn) error {
			return sd.Dialer.Control(network, address, c)
		}
	}
	if ctrl != nil {
		if err := ctrl(ctx, "unix", raddr.Name, denoRawConn{}); err != nil {
			return nil, err
		}
	}
	v, err := awaitJS(js.Global().Get("gotsNet").Call("connectUnix", raddr.Name))
	if err != nil {
		return nil, err
	}
	return newUnixConn(newDenoConn("unix", &denoSock{v: v})), nil
}

func (sl *sysListener) listenUnix(ctx context.Context, laddr *UnixAddr) (*UnixListener, error) {
	if sl.network != "unix" {
		return nil, os.NewSyscallError("socket", syscall.ENOTSUP)
	}
	if c := sl.ListenConfig.Control; c != nil {
		if err := c("unix", laddr.Name, denoRawConn{}); err != nil {
			return nil, err
		}
	}
	v := js.Global().Get("gotsNet").Call("listenUnix", laddr.Name)
	if v.Get("failed").Truthy() {
		return nil, denoError(v)
	}
	fd := &netFD{family: syscall.AF_UNIX, sotype: syscall.SOCK_STREAM, net: "unix", laddr: laddr, lis: &denoListener{v: v}}
	return &UnixListener{fd: fd, path: laddr.Name, unlink: true}, nil
}

// close removes the socket's file, which the system does not, and closes.
func (ln *UnixListener) close() error {
	ln.unlinkOnce.Do(func() {
		if ln.path != "" && ln.path[0] != '@' && ln.unlink {
			os.Remove(ln.path)
		}
	})
	return ln.fd.Close()
}

// unixAddr is the address the bridge gives as [path, -1].
func unixAddr(a js.Value) *UnixAddr {
	return &UnixAddr{Name: a.Index(0).String(), Net: "unix"}
}
