// Package water makes water for a pre-rendered isometric game, all in
// code: a normal map of small waves that tiles and loops, the shading of
// a water pixel as the game's camera sees it, with the sky mirrored by
// Fresnel's law, the sun glinting and foam at the shore, and the tiles of
// a shore, one for every mask terrain makes, water inside and the
// neighbouring land's edge given by the mask, over a ground tile.
package water
