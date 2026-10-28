//go:build gots_overlay

// Overlay for what net/http reaches by linkname: the limited MIME header
// reader, and serverHandler.ServeHTTP under another name.

package http

import "net/textproto"

func readMIMEHeader(r *textproto.Reader, maxMemory, maxHeaders int64) (textproto.MIMEHeader, error) {
	return textproto.GotsReadMIMEHeader(r, maxMemory, maxHeaders)
}

func badServeHTTP(sh serverHandler, rw ResponseWriter, req *Request) { sh.ServeHTTP(rw, req) }
