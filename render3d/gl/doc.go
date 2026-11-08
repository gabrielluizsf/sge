// Package gl is a render3d.Device drawn by OpenGL 4.2, core profile,
// through EGL: into an antui window on X11, or offscreen. It calls libEGL
// through package native, with cgo or without, on Linux and amd64;
// elsewhere Offscreen and Window say it cannot.
//
// A GL context belongs to one thread, so a device keeps one, locked, and
// does all its GL there: what a pass draws is recorded as it is set and
// handed over when the pass ends; the rest waits for its answer.
//
// render3d's clip space is WebGPU's, depth from 0 to 1, which GL 4.2
// cannot be told (it has no glClipControl): each vertex shader moves its
// depth into GL's, -1 to 1, so that the depth stored is render3d's. What
// each pipeline kind draws is render3d/soft's, to the pixel but for the
// edges of triangles; textures are sampled at their nearest pixel up
// close, as soft's are, and from their mipmaps far away.
package gl
