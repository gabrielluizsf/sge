//go:build gots_overlay

// Overlay for the method net/url exposes to the runtime by linkname.

package url

func badSetPath(u *URL, p string) error { return u.setPath(p) }
