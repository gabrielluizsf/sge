//go:build gots_overlay

// Overlay for the limited MIME header reader multipart links to.

package multipart

import "net/textproto"

func readMIMEHeader(r *textproto.Reader, maxMemory, maxHeaders int64) (textproto.MIMEHeader, error) {
	return textproto.GotsReadMIMEHeader(r, maxMemory, maxHeaders)
}
