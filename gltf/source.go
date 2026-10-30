package gltf

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// source is where a model's files are: the directory it may read them in,
// or none, for a model read from bytes, which may name only data: URIs.
type source struct {
	root   *os.Root
	limits Limits
}

// read is the file uri names: a data: URI's bytes, or a file under the
// model's directory, as os.Root opens it, so no "..", absolute path or
// link leads out of it. Nothing in it is expanded or run.
func (s source) read(uri string) ([]byte, error) {
	if rest, ok := strings.CutPrefix(uri, "data:"); ok {
		meta, payload, ok := strings.Cut(rest, ",")
		if !ok || !strings.HasSuffix(meta, ";base64") {
			return nil, malformed("data URI is not base64")
		}
		if int64(base64.StdEncoding.DecodedLen(len(payload))) > s.limits.Bytes {
			return nil, malformed("data URI over the limit")
		}
		return base64.StdEncoding.DecodeString(payload)
	}
	if s.root == nil {
		return nil, malformed("file %q named by a model read from bytes", uri)
	}
	name, err := localName(uri)
	if err != nil {
		return nil, err
	}
	f, err := s.root.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLimited(f, s.limits.Bytes)
}

// localName is uri as a path under the model's directory: percent-escapes
// undone, and refused unless it is local, with no scheme.
func localName(uri string) (string, error) {
	if strings.Contains(uri, ":") || strings.HasPrefix(uri, "/") || strings.Contains(uri, `\`) {
		return "", malformed("file %q is not a relative path", uri)
	}
	name := unescape(uri)
	if !filepath.IsLocal(name) {
		return "", malformed("file %q leaves the model's directory", uri)
	}
	return name, nil
}

// unescape undoes a URI's percent-escapes, leaving a broken one as it is.
func unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) && isHex(s[i+1]) && isHex(s[i+2]) {
			b.WriteByte(hexVal(s[i+1])<<4 | hexVal(s[i+2]))
			i += 2
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func isHex(c byte) bool { return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' }

func hexVal(c byte) byte {
	switch {
	case c <= '9':
		return c - '0'
	case c <= 'F':
		return c - 'A' + 10
	}
	return c - 'a' + 10
}

// readLimited reads r whole, failing past limit bytes.
func readLimited(r io.Reader, limit int64) ([]byte, error) {
	var b bytes.Buffer
	n, err := io.Copy(&b, io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if n > limit {
		return nil, malformed("file over %d bytes", limit)
	}
	return b.Bytes(), nil
}
