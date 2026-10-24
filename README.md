# sge

A small 2D game engine, written in Go on top of [antui](https://github.com/gabrielluizsf/antui),
that runs in two languages: as a Go module, and — transpiled to TypeScript by
the `gots` transpiler in this repository — as a Deno package on JSR.

The same engine, and the same game, behave the same in both: the
[`collect`](examples/collect) example plays frame for frame identically in Go
and on Deno, and the engine's Go tests pass on Deno too.

## Packages

| Package | What it does |
|---|---|
| `vector` | `Vector3` and `Quaternion` math |
| `transform` | position, rotation and scale |
| `collider` | box and circle overlap tests |
| `clock` | time between frames |
| `sprite` | antui canvases drawn with a pivot, scale and rotation; sprite sheets; PNG loading |
| `animation` | frame animations, played by name |
| `game` | the world: objects, scripts, collisions, drawing and the window loop |

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
nor `--allow-ffi`. See [`examples/ts/collect.ts`](examples/ts/collect.ts) for a
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
`Zero`, `MakeSlice`, `MakeMap` and `Convert`. `encoding/json` (with
`json/v2` underneath) and `encoding/binary` run on it. On Deno, JSON may nest
1000 levels deep rather than Go's 10000: a JavaScript stack does not grow,
and deeper input from outside is refused with an error instead of ending the
program.

## Limitations

- `reflect` cannot call functions or methods (`Value.Call`, `Method`), so
  `text/template`, `html/template`, `encoding/xml` and `encoding/gob` are not
  supported yet.
- A nil pointer to a struct, put in an interface, is a nil interface: a
  pointer to a struct is the object itself, and nil is nil.
- `net` and `net/http` are not supported: Deno's sockets are not wired in.
- `int` and `int64` are JavaScript numbers, exact up to 2^53; `uint64` is a
  `bigint`. Arithmetic with constants beyond 2^53, which `time` does on
  every date, is worked out exactly.
- Only Deno opens windows.

## License

MIT
