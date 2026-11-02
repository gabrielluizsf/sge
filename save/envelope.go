package save

import "encoding/json"

// envelope is what a file holds: the version of its format, and the data.
type envelope struct {
	Version int             `json:"version"`
	Data    json.RawMessage `json:"data"`
}
