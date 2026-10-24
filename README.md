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

Drawing, text, SVG and windows come from antui; `sge` does not draw anything
itself.

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

[`examples/farm`](examples/farm) is a farming game on all of it: buy seeds at
the stall, sow, water, pick and sell, on a farm in endless land. Every picture
in it is SVG, drawn from a fixed palette by
[`cmd/farmart`](examples/farm/cmd/farmart), HUD included; `go generate` in
`examples/farm` writes the sprites and a preview of them.

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

## How it works

`gots` (in [`gots/`](gots), run by [`cmd/gots`](cmd/gots)) type-checks Go
packages for `js/wasm` and writes TypeScript modules against a small runtime
for what JavaScript lacks: slices, pointers, Go's integer widths, interfaces,
panics and defers. Generics are instantiated per use. The standard library is
transpiled from Go's own source, except for what is native — `runtime`, `os`,
`sync`, `syscall/js` and a few more are written in TypeScript — and Go
overlays replace the handful of functions that need assembly or the runtime.

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

What blocks is worked out from the calls gots can see. A function value or
interface method that blocks is awaited when it is called from a function
that may block; called from one that never does, it is not waited for. A
`sync.Mutex` does not wait either: locking one that another goroutine holds
while it waits panics, rather than deadlock.

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

## Limitations

- `encoding/gob` is not supported yet. `reflect` cannot call a function
  that blocks, and a range over a function cannot have a body that blocks.
- A nil pointer to a struct, put in an interface, is a nil interface: a
  pointer to a struct is the object itself, and nil is nil.
- `net` and `net/http` are not supported: Deno's sockets are not wired in.
- `int` and `int64` are JavaScript numbers, exact up to 2^53; `uint64` is a
  `bigint`. Arithmetic with constants beyond 2^53, which `time` does on
  every date, is worked out exactly.
- Only Deno opens windows.

## License

MIT
