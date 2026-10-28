//go:build gots_overlay

package net

// ctrlNetwork is the network a Control hook is told, as Go names it by the
// address family: proto+"4" for an IPv4 address or none, proto+"6" else.
func ctrlNetwork(proto string, ip IP) string {
	if ip != nil && ip.To4() == nil {
		return proto + "6"
	}
	return proto + "4"
}
