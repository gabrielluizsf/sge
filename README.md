# sge

A small 2D game engine, written in Go on top of [antui](https://github.com/gabrielluizsf/antui),
that runs in two languages: as a Go module, and — transpiled to TypeScript by
the `gots` transpiler in this repository — as a Deno package on JSR.

The same engine, and the same game, behave the same in both: the
[`collect`](examples/collect) and [`farm`](examples/farm) examples play frame
for frame identically in Go and on Deno, and the engine's Go tests pass on
Deno too.

## Packages

| Package | What it does |
|---|---|
| `vector` | `Vector3` and `Quaternion` math |
| `transform` | position, rotation and scale |
| `collider` | box and circle overlap tests, and shapes of your own |
| `clock` | time between frames |
| `sprite` | antui canvases drawn with a pivot, scale and rotation; sprite sheets; PNG and SVG loading |
| `animation` | frame animations, played by name |
| `game` | the world: objects, scripts, collisions, the camera, drawing and the window loop |
| `tilemap` | grids of tiles: animated tiles, solid tiles, walking through layers of them |
| `terrain` | endless tile maps of named areas over a ground, made in chunks around the camera |
| `iso` | an isometric view of a flat world: cells drawn as diamonds, objects nearer over farther |
| `sound` | a mixer of clips and fading loops, and sounds made up from noise, tones, filters and envelopes |
| `scene` | switching a game between scenes, the world outside and a room, each kept as it was left |
| `daytime` | a clock whose hours go round, day and night light, a shade darkening the screen at night, and schedules by the hour |
| `town` | towns laid out over an endless world from a seed: streets, plazas, stalls, a smithy and houses of several kinds, woken near the camera |
| `pathfind` | shortest ways across a grid (A*), around walls without cutting corners, and steering along them |
| `inventory` | bags of items by name: put in, count, take out, a list at once or not at all |
| `request` | orders of random items drawn from a book of offers, for a reward, filled from a bag |
| `random` | a small random sequence that plays back the same in Go and on Deno |
| `label` | text in a TrueType face, smooth at its edges and outlined, as a picture for a HUD |
| `gltf` | glTF 2.0 and GLB models read as untrusted: meshes, materials, textures, nodes, skins, animations |
| `bake` | 3D models rendered ahead of time into an isometric game's sprites, lit by a sun casting shadows |

Drawing, text, SVG, windows and the speakers come from antui; `sge` does not
draw anything itself.

## Go

```go
g := game.New()
look, _ := sprite.Rect(24, 24, canvas.RGB(80, 170, 255))
hero := g.Add(game.NewObject("player"))
hero.Sprite = look
hero.AddScript(game.ScriptFunc(func(o *game.Object, dt float64) {
	if o.Game().Window().KeyDown(antui.KeyRight) {
		o.Transform.Translate(vector.Vector3{X: 200 * dt})
	}
}))
g.Run("my game", 480, 320)
```

A script is any value with `Update(o *Object, dt float64)`; it may also have
`Start(o)`, `Collision(o, other)` and `Paint(o, cv)`.

`Game.Camera` is the point of the world at the window's top-left corner;
`game.Follow` keeps an object in the middle of the window, whatever its size,
and objects with `Screen` set, a HUD, stay put. Painters draw at points of the
screen: `Game.ToScreen` says where a point of the world is.

## Tile maps and endless terrain

A `tilemap.Map` is a `game.Object` whose sprite is its tiles: build one from a
tile set and text, animate tiles (grass that sways, cell by cell, or water that
ripples all together), mark tiles solid, and walk an object through several
layers at once with `tilemap.Move`. A map of opaque tiles, the ground, is
copied onto the window with antui's `Copy` instead of blended; a layer of
things standing on it draws only its tiles.

`terrain` makes maps that never end. The ground covers everything, and named
areas lie over it, each on its layer, worked out for any cell from a seed:

```go
gen := &terrain.Generator{Tiles: tiles, Layers: 3}
gen.Add(terrain.Area{Name: "grass", Layer: 0, In: terrain.Everywhere, Tile: grassTile})
lakes := terrain.Open(terrain.Blob(seed, 9, 0.6)) // round, nothing thinner than 3 cells
gen.Add(terrain.Area{Name: "water", Layer: 1, In: lakes, Tile: terrain.Edged(lakes, shoreTile)})
trees := terrain.Scatter(seed+1, 0.05, 2)
gen.Add(terrain.Area{Name: "trees", Layer: 2, Tile: treeTile,
	In: func(c, r int) bool { return trees(c, r) && !gen.Near("water", c, r, 1) }})
world := terrain.New("world", gen, 16, 0, 1, 2)
g.Add(world.Object)                  // after what moves the camera
g.Covered = true                     // the world fills the window: no clearing
hero.AddScript(game.Follow{Lag: 0.1})
tilemap.Move(hero, dx, dy, world)    // its water is a wall
```

The world makes the chunks the window shows at once and those around it a few
a frame, and lets go of those left behind, so a window of any size, maximised
or resized, is always filled. `Edged` picks a cell's tile from its neighbours,
one of 47, so shores join up; `Smooth`, `Open` and `Memo` shape areas and keep
their answers while a chunk is made.

## Isometric worlds

`iso` sees a world from a corner. The world stays a flat grid, where objects
move and collide as in any game; only its picture turns, each square cell a
diamond twice as wide as it is tall:

```go
g.Projection = iso.New(32, 64, 32) // cells 32 square, drawn 64 by 32
tiles.CellWidth, tiles.CellHeight = 32, 32 // tile pictures 64 wide, taller where things stand
trees.Sorted = true                // each tree among the objects of its Z
hero.Z = trees.Z                   // drawn nearer over farther, by where it stands
cell := g.ToWorld(mouse)           // the ground under the mouse
```

A projected map stands each tile on the near corner of its cell's diamond, so
a tile taller than its diamond, a tree or a fence, reaches up out of it. A
`Sorted` map's tiles are each drawn among the objects of its Z by the screen
height of where they stand, as objects are by their `Footing`: a farmer in
front of a trunk is over it, one behind it under it. `Follow` centres the
camera, `VisibleWorld` is the patch of ground the view shows and `terrain`
makes the chunks under it. Ground tiles, a run of opaque pixels on every row,
are copied rather than blended, and every tile's air is cut away once. A
`tilemap.NewLazyTileSet` draws each tile the first time a map shows it, so a
game of many tiles, drawn from SVG, starts at once.

## Towns, people and time

`town` lays out towns over an endless world: a region of cells each, with a
town or none as the seed says, and a `Keeper` script that wakes a town's
objects as the camera nears it and lets them go when it leaves. `pathfind`
walks people round what stands in the way, `daytime` gives them hours to keep
and darkens the night, `inventory` and `request` are what they carry and what
they ask for, and `scene` takes the player into a house and back out to the
world as it was. A `terrain.World` can be edited for good (`World.Set`, a tree
felled stays felled), and stand objects on its chunks (`World.Props`).
A `sprite.Sprite` with `Make` is drawn the first time it is shown, so a
character's many frames cost nothing until seen.

[`examples/farm`](examples/farm) is a farming game on all of it, seen from the
corner: buy seeds at the stall, sow, water, pick and sell, on a farm in endless
land. Along its path lie towns of women and men of every skin, each with a
trade they keep by the hour: farmers tend their gardens and sell at the stall,
lumberjacks fell trees and sell their wood, the smith forges at the anvil,
feeding the forge wood they cut or buy, and others stroll about town; at
night everyone sleeps at home. Villagers ask for wood and crops for a reward,
the lumberjack trades wood, and the smith gives a stone axe for an errand and
sells iron, gold and diamond ones, each felling a tree in fewer strokes. Every
house can be gone into. The player sleeps in their own bed; one still up a
minute past midnight faints and loses all their money. Arrows walk the way
they point on the screen; space, or a click on a cell within reach, which the
mouse outlines, uses what the farmer holds; E talks, goes in, sleeps or
shops. Every
picture in it is SVG, drawn from a fixed palette by
[`cmd/farmart`](examples/farm/cmd/farmart), HUD included, the ground laid on
the diamonds and what stands on it drawn upright; `go generate` in
`examples/farm` writes the sprites and a preview of them.

## Sound

A game with a `sound.Mixer` plays it while it runs: clips once, each with
its gain, pan and pitch, and loops that fade in and out.

```go
g.Sound = sound.NewMixer(sound.Rate)
step := sound.NewClip(sound.Noise(1, sound.Seconds(0.1)), 0.5) // or made up from tones, filters and envelopes
g.Sound.Play(step, sound.Voice{Gain: 0.4, Pan: -0.3, Pitch: 1.05})
wind := g.Sound.Loop(windClip)
wind.Fade(0.2) // eases in; Fade(0) eases out
```

The mix keeps every sample from -1 to 1, so no sound can click the speaker,
and where there are no speakers the game runs in silence. Natively antui plays
it through PulseAudio, WASAPI, CoreAudio or AAudio; on Deno through SDL2's
audio queue, with `--allow-ffi`; in a browser with Web Audio. The farm's
steps, water, leaves, birds and HUD sounds are all made up in
[`examples/farm/sfx`](examples/farm/sfx).

## TypeScript (Deno)

```ts
import { antui, canvas, game, sprite, vector } from "jsr:@gabrielluizsf/sge";

const g = game.New();
const [look] = sprite.Rect(24, 24, canvas.RGB(80, 170, 255));
const hero = g.Add(game.NewObject("player"));
hero.Sprite = look;
hero.AddScript({
  Update(o: game.Object, dt: number) {
    if (o.Game()?.Window()?.KeyDown(antui.KeyRight)) o.Transform.Translate(new vector.Vector3(200 * dt));
  },
});
await g.Run("my game", 480, 320);
```

Run it with `deno run --allow-ffi --allow-env --allow-read`. Windows open
through SDL2, which has to be installed; drawing off screen needs neither SDL2
nor `--allow-ffi`. See [`examples/collect/collect.ts`](examples/collect/collect.ts) for a
whole game.

Grant no more than the game needs. `--allow-ffi` lets native code run outside
Deno's sandbox, so give it only to games you trust, and scope the rest:
`--allow-read=./assets` rather than `--allow-read`, and `--allow-net` only
with the hosts an SVG may load pictures from (`--allow-net=cdn.example.com`),
since antui's SVG renderer follows the `href` of the pictures it draws (an
SVG that came from the network may not name local files, but one read from
disk may name both). Never use `-A` for SVGs or other assets you did not
write.

The TypeScript API is the Go API, transpiled: the same names, structs as
classes, multiple results as tuples, `nil` as `null`. Text crosses as plain
JavaScript strings — function arguments and results, string fields, and the
strings a TypeScript script hands back to Go — and `goSlice` makes a Go slice
of an array.

## Browser

A Go program runs in a browser as it runs on Deno: its window is a
`<canvas>` on the page and its input the page's keys, mouse and wheel.

```sh
go run ./cmd/gots -web site ./examples/farm
# site/app.js and site/index.html: serve the directory, open the page
```

`-web` bundles the program into `app.js` with `deno bundle` and writes an
`index.html` beside it, unless the directory has one already: a page of
your own keeps working, with a `<canvas data-gots>` where the game should
draw (one is added to the body otherwise). `fmt` prints to the console.
A window's icon is the page's, its cursor the canvas's, notifications the
browser's, and the system's dark theme and the display's scale what the
page is told.

What a browser keeps from pages is the browser's to keep: sound starts with
the player's first click or key, full screen is granted only in answer to
one, notifications only once the
player allows them, the clipboard is what the player last pasted on the
canvas, dropped files have no paths to open, and `net` has no sockets.

## How it works

`gots` (in [`gots/`](gots), run by [`cmd/gots`](cmd/gots)) type-checks Go
packages for `js/wasm` and writes TypeScript modules against a small runtime
for what JavaScript lacks: slices, pointers, Go's integer widths, interfaces,
panics and defers. Generics are instantiated per use. The standard library is
transpiled from Go's own source, except for what is native — `runtime`, `os`,
`sync`, `syscall/js` and a few more are written in TypeScript — and Go
overlays replace the handful of functions that need assembly or the runtime.
Go's doc comments become the modules' JSDoc, the package's as the module's,
so the TypeScript is documented as the Go is; what only the transpiled code
uses is marked `@ignore`.

`go generate .` writes the package to `deno/`. `go test ./...` runs, besides
the Go tests, conformance programs that must print the same under `go run` and
Deno, and the engine's tests transpiled to Deno.

## Goroutines and channels

Goroutines, channels, `select`, `sync.WaitGroup`, `sync.Cond` and the
timers of `time` work as in Go, on JavaScript's one thread: a function that
may block (a send, a receive, a `select` that waits, `time.Sleep`,
`WaitGroup.Wait`, or a call of a function that blocks) becomes an `async`
function its callers await, and everything else stays as fast as plain
JavaScript. A goroutine is a microtask; a program ends when `main` returns,
dies of a panic no goroutine recovers, and reports a deadlock as Go does.

From TypeScript, the functions that block return promises:
`await g.Run(...)`, and the game loop lets Deno's event loop run between
frames.

What blocks is worked out from the calls gots can see. A function that
blocks only through what it calls by interface or function value, say
`bufio.Reader.Read`, which blocks when it reads a socket and not when it
reads a string, is written twice: as itself, which never waits, for code
that cannot, and as its twin, which awaits what it calls, for code that can
(a goroutine, `main`, a test, an HTTP handler). Engine code drawing a frame
stays synchronous; a server's goroutine waits on its socket. Should code
that cannot wait reach something that does block, it panics saying so,
rather than carry on with a promise for a result. `sync.Mutex` and
`sync.RWMutex` wait their turn in code that can wait, with no turn of the
event loop when nobody holds them; code that cannot wait, which only a
lock held across a wait elsewhere would make wait, panics saying so.

## Networking

`net` dials and listens over TCP, Unix stream and datagram sockets and
UDP with Deno's sockets, and `net/http` serves with Go's own server on them
and sends with Go's own `Transport`: proxies, HTTP/2, `DialContext`,
everything as Go has it. Host names go to Deno unresolved, so
`--allow-net` checks the name and the addresses it resolves to: grant only
the hosts a game talks to (`--allow-net=127.0.0.1:8080,api.example.com`).
A `Dialer` with a `Control` hook resolves the name, asks the hook about
each address and dials the one it lets through; one with a `Resolver` of
its own resolves with it, Go's DNS client over Deno's sockets. Reads and
writes hold their deadlines, and a read or write that runs out of time
loses nothing of the stream. A `Content-Length` past 2^53, which a
JavaScript number would round, is refused by the server and the client
rather than read as another length.

`crypto/tls` is Go's own, whole: versions, suites, curves (X25519MLKEM768
included), session resumption, client certificates, OCSP stapling,
`ExportKeyingMaterial`, `VerifyPeerCertificate` and `VerifyConnection`
all behave as in Go. `x509.SystemCertPool()` holds the roots Deno itself
verifies with. In a browser, which has no sockets, `http.Client` sends
through `fetch`, with the client, not `fetch`, following redirects so that
`CheckRedirect` and Go's rules for the headers that may follow a redirect
to another host hold; a `Transport` with a `TLSClientConfig` fails there,
rather than have `fetch` drop its roots and pins.

Cryptography runs at JavaScript's speed rather than at that of Go's 64-bit
words as BigInts: the NIST curves' and Curve25519's fields are unrolled
22-bit limbs, bigmod, GHASH, SHA-512 and Keccak are the runtime's on
32-bit words, and AES-GCM, ChaCha20-Poly1305 and X25519 are Deno's own
(`node:crypto`) where it has them. A P-256 signature checks in some 5ms,
Ed25519 in 3ms, and TLS moves tens of megabytes a second.

`encoding/gob`, `encoding/json` and `encoding/xml` refuse input nested
deeper than a JavaScript stack holds with an error, rather than end the
program: Go's limits are sized for a stack that grows.

## reflect

`reflect` works over the transpiler's run-time type descriptors: types,
kinds, struct fields and their tags, `Implements`, and values read and set
through pointers, fields, slices, arrays, maps and interfaces, `New`,
`Zero`, `MakeSlice`, `MakeMap` and `Convert`, and calls: `Value.Call`,
`Method`, `MethodByName`, and ranging with `Seq`. `encoding/json` (with
`json/v2` underneath), `encoding/xml`, `encoding/binary`, `text/template`
and `html/template` run on it; `html/template`'s escaping is checked to be
Go's byte for byte. On Deno, JSON may nest
1000 levels deep rather than Go's 10000, and XML 256: a JavaScript stack
does not grow, and deeper input from outside is refused with an error
instead of ending the program.

## Limits of Deno

gots itself has none of its own; what follows is Deno's.

- `os.Root` checks a name a step at a time and then opens it, as Deno has no
  `openat`: a link swapped in between the check and the open is not seen,
  which Go's `Root`, opening each step from the last, would see.
- Deno's permissions: `--allow-net` grants what a program connects to,
  by name and by address, and Unix sockets by path
  (`--allow-net=unix:/run/game.sock`); `--unstable-net` lets UDP and Unix
  datagram sockets in; `--allow-sys=ca` lets Go read the roots Deno
  verifies with, for TLS with the system's roots; `--allow-env` and
  `--allow-read` what `os` reads.
- Deno's sockets have no descriptors, so a `Control` hook that asks for one
  is told there is none; Unix packet sockets (`unixpacket`) do not exist
  in Deno; every Unix datagram socket is bound to a path, so one dialed
  unnamed is bound to a temporary one, which its peer sees where Go's
  would see none. A UDP socket dialing a host off this machine is bound to
  every address, `0.0.0.0`, as Deno binds before it sends. Closing a
  listener while an `Accept` waits can reset, rather than refuse, a
  connection dialed in that same moment: Deno closes the socket once the
  pending accept is done.
- JavaScript's numbers, which Deno computes with: `int` and `int64` are
  numbers, exact up to 2^53, as the cost of a `bigint` in every sum would
  slow down all of them; `uint64` is a `bigint`, for integers that need
  more, and so is `math/big`. `uint` is 32 bits, as on 386 and arm, which a
  number holds exactly, and so `math.MaxInt` is a 32-bit machine's,
  2^31-1, as Go works it out from `uint`.
- JavaScript's stack, which does not grow: JSON nests 1000 levels deep
  rather than Go's 10000, and XML 256, deeper input refused with an error.
- Deno has no windows of its own: they open through SDL2, which Deno loads
  by FFI (`--allow-ffi`), and SDL2 has no notifications or theme to tell
  of, so those come from the system's own tools (`notify-send` and
  `gsettings`, `osascript` and `defaults`, PowerShell and `reg`), which
  Deno runs with `--allow-run`.

## License

MIT
