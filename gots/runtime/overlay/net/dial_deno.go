//go:build gots_overlay

// Overlay for dialing and listening on TCP: Deno's sockets, through the
// gotsNet bridge, rather than the in-memory network Go keeps for js.
// Dialing hands Deno the host name as it is, so Deno looks it up and its
// --allow-net checks the name and the addresses it resolves to.

package net

import (
	"context"
	"internal/nettrace"
	"os"
	"strconv"
	"syscall"
	"syscall/js"
)

func (d *Dialer) DialContext(ctx context.Context, network, address string) (Conn, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return d.dialDeno(ctx, network, address)
	}
	ctx, cancel := d.dialCtx(ctx)
	defer cancel()

	// Shadow the nettrace (if any) during resolve so Connect events don't fire for DNS lookups.
	resolveCtx := ctx
	if trace, _ := ctx.Value(nettrace.TraceKey{}).(*nettrace.Trace); trace != nil {
		shadow := *trace
		shadow.ConnectStart = nil
		shadow.ConnectDone = nil
		resolveCtx = context.WithValue(resolveCtx, nettrace.TraceKey{}, &shadow)
	}

	addrs, err := d.resolver().resolveAddrList(resolveCtx, "dial", network, address, d.LocalAddr)
	if err != nil {
		return nil, &OpError{Op: "dial", Net: network, Source: nil, Addr: nil, Err: err}
	}

	sd := &sysDialer{
		Dialer:  *d,
		network: network,
		address: address,
	}

	var primaries, fallbacks addrList
	if d.dualStack() && network == "tcp" {
		primaries, fallbacks = addrs.partition(isIPv4)
	} else {
		primaries = addrs
	}

	return sd.dialParallel(ctx, primaries, fallbacks)
}

// dialDeno dials host:port with Deno.connect, giving up when ctx or the
// dialer's timeout says to.
func (d *Dialer) dialDeno(ctx context.Context, network, address string) (Conn, error) {
	fail := func(err error) (Conn, error) {
		return nil, &OpError{Op: "dial", Net: network, Err: err}
	}
	// Deno connects by name, past anything that would look at the address
	// first: a Control or ControlContext, which servers fetching URLs from
	// outside refuse private addresses in, and a Resolver of the program's
	// own. Rather than skip them, which would let through what they
	// refuse, a dial with one fails, as Go's own js port does.
	if d.Control != nil || d.ControlContext != nil || d.Resolver != nil && d.Resolver != DefaultResolver {
		return fail(os.NewSyscallError("socket", syscall.ENOTSUP))
	}
	host, service, err := SplitHostPort(address)
	if err != nil {
		return fail(err)
	}
	port, err := denoPort(network, service)
	if err != nil {
		return fail(err)
	}
	ctx, cancel := d.dialCtx(ctx)
	defer cancel()
	type result struct {
		v   js.Value
		err error
	}
	done := make(chan result, 1)
	go func() {
		v, err := awaitJS(js.Global().Get("gotsNet").Call("connect", host, port))
		done <- result{v, err}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			return fail(r.err)
		}
		fd := newDenoConn(network, &denoSock{v: r.v})
		return newTCPConn(fd, d.KeepAlive, d.KeepAliveConfig, testPreHookSetKeepAlive, testHookSetKeepAlive), nil
	case <-ctx.Done():
		go func() {
			// A connection made after all is closed.
			if r := <-done; r.err == nil {
				r.v.Call("close")
			}
		}()
		return fail(mapErr(ctx.Err()))
	}
}

// denoPort is the port service names: a number, or a well-known service.
func denoPort(network, service string) (int, error) {
	if p, needsLookup := parsePort(service); !needsLookup {
		if p < 0 || p > 0xFFFF {
			return 0, &AddrError{Err: "invalid port", Addr: service}
		}
		return p, nil
	}
	return lookupPortMap(network, service)
}

func (sd *sysDialer) doDialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
	c, err := sd.Dialer.dialDeno(ctx, sd.network, JoinHostPort(raddr.IP.String(), strconv.Itoa(raddr.Port)))
	if err != nil {
		return nil, err
	}
	return c.(*TCPConn), nil
}

func (sl *sysListener) listenTCP(ctx context.Context, laddr *TCPAddr) (*TCPListener, error) {
	if sl.ListenConfig.Control != nil {
		// Deno's listen has no socket to show it: fail, as Go's js port does.
		return nil, &OpError{Op: "listen", Net: sl.network, Addr: laddr, Err: os.NewSyscallError("socket", syscall.ENOTSUP)}
	}
	host := ""
	if laddr != nil && laddr.IP != nil && !laddr.IP.IsUnspecified() {
		host = laddr.IP.String()
	}
	port := 0
	if laddr != nil {
		port = laddr.Port
	}
	v, err := denoListen(host, port)
	if err != nil {
		return nil, &OpError{Op: "listen", Net: sl.network, Addr: laddr, Err: err}
	}
	l := &denoListener{v: v}
	fd := &netFD{family: syscall.AF_INET, sotype: syscall.SOCK_STREAM, net: sl.network, laddr: tcpAddr(v.Call("local")), lis: l}
	return &TCPListener{fd: fd, lc: sl.ListenConfig}, nil
}

// denoListen calls the bridge's listen, which gives an error rather than
// throw one.
func denoListen(host string, port int) (js.Value, error) {
	v := js.Global().Get("gotsNet").Call("listen", host, port)
	if v.Get("failed").Truthy() {
		return js.Value{}, denoError(v)
	}
	return v, nil
}
