//go:build gots_overlay

// Overlay exporting the MIME header reader with limits that mime/multipart
// and net/http reach by linkname, which transpiled packages cannot do: it
// is textproto's own, holding the memory and header count limits.

package textproto

// GotsReadMIMEHeader is readMIMEHeader.
func GotsReadMIMEHeader(r *Reader, maxMemory, maxHeaders int64) (MIMEHeader, error) {
	return readMIMEHeader(r, maxMemory, maxHeaders)
}
