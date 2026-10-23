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
g.Run("my game", 480, 320);
```

Run it with `deno run --allow-ffi --allow-env --allow-read`. Windows open
through SDL2, which has to be installed; drawing off screen needs neither SDL2
nor `--allow-ffi`. See [`examples/ts/collect.ts`](examples/ts/collect.ts) for a
whole game.

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

## Limitations

- Goroutines and channels are not supported yet.
- `reflect` is not supported, and so neither are `encoding/json` and the
  packages built on it.
- `int` and `int64` are JavaScript numbers, exact up to 2^53; `uint64` is a
  `bigint`.
- Only Deno opens windows; on Deno a game loop holds the event loop, so
  promises (such as antui fetching an SVG's remote pictures) settle once it
  ends.

## License

MIT
