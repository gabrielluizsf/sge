// Package save keeps a game's saved state on disk: JSON files, one per
// fixed name, in the game's own folder under the user's config directory.
// Writing is atomic, so a save that fails leaves the old one whole; reading
// takes no more than a limit of bytes, refuses unknown fields, trailing
// data and a version it does not expect, so a hostile file gives an error.
package save
