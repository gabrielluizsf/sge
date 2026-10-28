//go:build gots_overlay

// Overlay for the Content-Length a server reads: one past what a
// JavaScript number holds exactly is refused, rather than read as another
// length than the client meant, the start of request smuggling.

package http

import (
	"net/textproto"
	"strconv"
)

func parseContentLength(clHeaders []string) (int64, error) {
	if len(clHeaders) == 0 {
		return -1, nil
	}
	cl := textproto.TrimString(clHeaders[0])

	// The Content-Length must be a valid numeric value.
	// See: https://datatracker.ietf.org/doc/html/rfc2616/#section-14.13
	if cl == "" {
		if httplaxcontentlength.Value() == "1" {
			httplaxcontentlength.IncNonDefault()
			return -1, nil
		}
		return 0, badStringError("invalid empty Content-Length", cl)
	}
	n, err := strconv.ParseUint(cl, 10, 63)
	if err != nil {
		return 0, badStringError("bad Content-Length", cl)
	}
	if n > maxExactLength {
		return 0, badStringError("Content-Length too large", cl)
	}
	return int64(n), nil
}
