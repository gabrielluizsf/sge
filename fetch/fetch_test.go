package fetch

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sum(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

// server serves body at /model.bin, a redirect elsewhere at /away, and
// counts requests.
func server(t *testing.T, body []byte) (*httptest.Server, *int) {
	n := 0
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n++
		switch r.URL.Path {
		case "/model.bin":
			w.Write(body)
		case "/away":
			http.Redirect(w, r, "https://elsewhere.example/x", http.StatusFound)
		case "/plain":
			http.Redirect(w, r, "http://"+r.Host+"/model.bin", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &n
}

func fetcher(t *testing.T, srv *httptest.Server) Fetcher {
	return Fetcher{Cache: t.TempDir(), Hosts: []string{strings.TrimPrefix(srv.URL, "https://")}, Max: 1 << 20, Client: srv.Client()}
}

// A file is fetched into the cache, and taken from there next time.
func TestAFileIsFetchedOnceAndKept(t *testing.T) {
	body := []byte("vertices")
	srv, n := server(t, body)
	fe := fetcher(t, srv)
	f := File{URL: srv.URL + "/model.bin", Path: "tree/model.bin", SHA256: sum(body)}
	for range 2 {
		p, err := fe.Get(f)
		if err != nil {
			t.Fatal(err)
		}
		if got, _ := os.ReadFile(p); string(got) != "vertices" {
			t.Fatalf("kept %q", got)
		}
	}
	if *n != 1 {
		t.Errorf("fetched %d times", *n)
	}
}

// Bytes other than those pinned are refused and not kept.
func TestAFileThatIsNotWhatWasPinnedIsRefused(t *testing.T) {
	srv, _ := server(t, []byte("tampered"))
	fe := fetcher(t, srv)
	f := File{URL: srv.URL + "/model.bin", Path: "model.bin", SHA256: sum([]byte("original"))}
	if _, err := fe.Get(f); !errors.Is(err, ErrMismatch) {
		t.Fatalf("got %v", err)
	}
	if _, err := os.Stat(filepath.Join(fe.Cache, "model.bin")); err == nil {
		t.Error("the tampered file was kept")
	}
}

// Only HTTPS from a trusted host is fetched, redirects too, nothing over
// the limit, and nothing kept outside the cache.
func TestOnlyTrustedFetchesAreMade(t *testing.T) {
	body := []byte("vertices")
	srv, _ := server(t, body)
	fe := fetcher(t, srv)
	host := strings.TrimPrefix(srv.URL, "https://")
	for name, f := range map[string]File{
		"plain HTTP":          {URL: "http://" + host + "/model.bin", Path: "a", SHA256: sum(body)},
		"untrusted host":      {URL: "https://example.com/model.bin", Path: "a", SHA256: sum(body)},
		"redirect elsewhere":  {URL: srv.URL + "/away", Path: "a", SHA256: sum(body)},
		"redirect to HTTP":    {URL: srv.URL + "/plain", Path: "a", SHA256: sum(body)},
		"path out of cache":   {URL: srv.URL + "/model.bin", Path: "../escape.bin", SHA256: sum(body)},
		"absolute cache path": {URL: srv.URL + "/model.bin", Path: "/tmp/escape.bin", SHA256: sum(body)},
		"credentials in URL":  {URL: "https://user:pw@" + host + "/model.bin", Path: "a", SHA256: sum(body)},
	} {
		if _, err := fe.Get(f); err == nil {
			t.Errorf("%s: fetched", name)
		}
	}
	fe.Max = 4
	if _, err := fe.Get(File{URL: srv.URL + "/model.bin", Path: "big", SHA256: sum(body)}); err == nil {
		t.Error("a file over the limit was fetched")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(fe.Cache), "escape.bin")); err == nil {
		t.Error("a file was written outside the cache")
	}
}
