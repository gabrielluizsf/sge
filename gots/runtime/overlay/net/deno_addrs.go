//go:build gots_overlay

// The addresses a Deno dial tries: see dialDeno.

package net

import (
	"context"
	"syscall/js"
)

// denoAddrs is what host resolves to, by the program's own Resolver if
// own is true and by Deno otherwise, or host itself if it is an address.
func (d *Dialer) denoAddrs(ctx context.Context, host string, own bool) ([]string, error) {
	if ParseIP(host) != nil {
		return []string{host}, nil
	}
	if own {
		addrs, err := d.Resolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		ips := make([]string, len(addrs))
		for i, a := range addrs {
			ips[i] = a.IP.String()
		}
		return ips, nil
	}
	v, err := awaitJS(js.Global().Get("gotsNet").Call("lookup", host))
	if err != nil {
		return nil, err
	}
	ips := make([]string, v.Length())
	for i := range ips {
		ips[i] = v.Index(i).String()
	}
	return ips, nil
}
