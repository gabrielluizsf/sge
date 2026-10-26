//go:build gots_overlay

// Overlay for UDP: Deno's datagram sockets, through the gotsNet bridge.
// Deno has them behind --unstable-net; without it listening or dialing
// UDP fails saying so.

package net

import (
	"context"
	"os"
	"syscall"
	"syscall/js"
)

func (sl *sysListener) listenUDP(ctx context.Context, laddr *UDPAddr) (*UDPConn, error) {
	if c := sl.ListenConfig.Control; c != nil {
		if err := c(ctrlNetwork("udp", laddr.IP), laddr.String(), denoRawConn{}); err != nil {
			return nil, err
		}
	}
	return openUDP(sl.network, laddr, nil)
}

func (sd *sysDialer) dialUDP(ctx context.Context, laddr, raddr *UDPAddr) (*UDPConn, error) {
	ctrl := sd.Dialer.ControlContext
	if ctrl == nil && sd.Dialer.Control != nil {
		ctrl = func(ctx context.Context, network, address string, c syscall.RawConn) error {
			return sd.Dialer.Control(network, address, c)
		}
	}
	if ctrl != nil {
		if err := ctrl(ctx, ctrlNetwork("udp", raddr.IP), raddr.String(), denoRawConn{}); err != nil {
			return nil, err
		}
	}
	if laddr == nil {
		// Bound where the peer is reached from: loopback for loopback, which
		// --allow-net grants with the peer, every address otherwise.
		laddr = &UDPAddr{}
		switch {
		case raddr.IP.IsLoopback() && raddr.IP.To4() != nil:
			laddr.IP = IPv4(127, 0, 0, 1)
		case raddr.IP.IsLoopback():
			laddr.IP = IPv6loopback
		case raddr.IP.To4() == nil:
			laddr.IP = IPv6unspecified
		}
	}
	return openUDP(sd.network, laddr, raddr)
}

// openUDP opens Deno's UDP socket on laddr, connected to raddr if given.
func openUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error) {
	host := ""
	if laddr != nil && laddr.IP != nil && !laddr.IP.IsUnspecified() {
		host = laddr.IP.String()
	} else if laddr != nil && laddr.IP != nil && laddr.IP.To4() == nil {
		host = "::"
	}
	port := 0
	if laddr != nil {
		port = laddr.Port
	}
	v := js.Global().Get("gotsNet").Call("listenUdp", host, port)
	if v.Get("failed").Truthy() {
		if v.Get("name").String() == "NotSupported" {
			return nil, os.NewSyscallError("socket", syscall.ENOTSUP)
		}
		return nil, denoError(v)
	}
	l := v.Call("local")
	local := &UDPAddr{IP: ParseIP(l.Index(0).String()), Port: l.Index(1).Int()}
	family := syscall.AF_INET
	if local.IP.To4() == nil {
		family = syscall.AF_INET6
	}
	d := &denoDgram{v: v, peer: raddr}
	fd := &netFD{family: family, sotype: syscall.SOCK_DGRAM, net: network, laddr: local, isConnected: raddr != nil, dgram: d}
	if raddr != nil {
		fd.raddr = raddr
	}
	return newUDPConn(fd), nil
}
