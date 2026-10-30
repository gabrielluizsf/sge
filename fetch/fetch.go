// Package fetch downloads the files a game's assets are made from, as a
// build step: over HTTPS from hosts it is told to trust, redirects too,
// no bigger than a limit, each checked against the SHA-256 its manifest
// pins before it is kept, in a cache outside the repository, each where
// its manifest says, which may not lead out of the cache.
package fetch

import "net/http"

// File is a file to fetch: where from, where in the cache it is kept,
// and the SHA-256 its bytes must have, in hex.
type File struct {
	URL    string `json:"url"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// Fetcher fetches into Cache, a directory, from Hosts only, by name, or
// name:port for a port other than 443, no file over Max bytes, with Client (http.DefaultClient when nil).
type Fetcher struct {
	Cache  string
	Hosts  []string
	Max    int64
	Client *http.Client
}
