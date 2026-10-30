package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"slices"
)

// Allowed reports, as an error, whether raw may be fetched: HTTPS, from one
// of the hosts trusted, which names a port other than HTTPS's own when
// it may be fetched from one.
func (fe Fetcher) Allowed(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "https" || u.User != nil || !slices.Contains(fe.Hosts, u.Host) {
		return fmt.Errorf("fetch: %s is not an HTTPS URL of a host trusted", raw)
	}
	return nil
}

// download is the body at raw, refused past Max bytes.
func (fe Fetcher) download(raw string) ([]byte, error) {
	if err := fe.Allowed(raw); err != nil {
		return nil, err
	}
	resp, err := fe.client().Get(raw)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetch: %s: %s", raw, resp.Status)
	}
	if resp.ContentLength > fe.Max {
		return nil, fmt.Errorf("fetch: %s is %d bytes, over %d", raw, resp.ContentLength, fe.Max)
	}
	var b bytes.Buffer
	n, err := io.Copy(&b, io.LimitReader(resp.Body, fe.Max+1))
	if err != nil {
		return nil, err
	}
	if n > fe.Max {
		return nil, fmt.Errorf("fetch: %s is over %d bytes", raw, fe.Max)
	}
	return b.Bytes(), nil
}
