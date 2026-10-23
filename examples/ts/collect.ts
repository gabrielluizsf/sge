// collect, written in TypeScript against sge's transpiled packages: the
// same game as examples/collect, scripts and all, in the other language.
//
//   deno run --allow-ffi --allow-env --allow-read examples/ts/collect.ts
import { antui, canvas, collider, game, goString, sprite, vector } from "../../deno/mod.ts";

const width = 480, height = 320;

// Player is a script: any class with an Update method is one, and the
// optional Collision and Paint methods are picked up the same way.
class Player {
  score = 0;
  constructor(readonly speed: number) {}

  Update(o: game.Object, dt: number): void {
    const win = o.Game()?.Window();
    if (!win) return;
    let dir = new vector.Vector3();
    if (win.KeyDown(antui.KeyLeft) || win.KeyDown(antui.KeyA)) dir = dir.Add(new vector.Vector3(-1));
    if (win.KeyDown(antui.KeyRight) || win.KeyDown(antui.KeyD)) dir = dir.Add(new vector.Vector3(1));
    if (win.KeyDown(antui.KeyUp) || win.KeyDown(antui.KeyW)) dir = dir.Add(new vector.Vector3(0, -1));
    if (win.KeyDown(antui.KeyDown) || win.KeyDown(antui.KeyS)) dir = dir.Add(new vector.Vector3(0, 1));
    if (dir.SqrLength() > 0) o.Transform.Translate(dir.Normalized().Scale(this.speed * dt));
    const p = o.Transform.Position;
    p.X = Math.min(Math.max(p.X, 16), width - 16);
    p.Y = Math.min(Math.max(p.Y, 16), height - 16);
  }

  Collision(_o: game.Object, other: game.Object): void {
    if (other.Name === "coin") {
      this.score++;
      place(other);
    }
  }

  Paint(_o: game.Object, cv: canvas.Canvas): void {
    cv.Text(new canvas.Point(8, 8), goString(`pontos ${this.score}`), canvas.RGB(240, 240, 240));
  }
}

function place(coin: game.Object): void {
  coin.Transform.Position = new vector.Vector3(20 + Math.random() * (width - 40), 30 + Math.random() * (height - 50));
}

const g = game.New();
g.Background = canvas.RGB(20, 24, 32);

const [body, err] = sprite.Rect(24, 24, canvas.RGB(80, 170, 255));
if (err !== null) throw new Error(err.Error());
const [coinLook] = sprite.Rect(12, 12, canvas.RGB(250, 200, 40));

const player = new Player(220);
const hero = g.Add(game.NewObject("player"))!;
hero.Sprite = body;
hero.Z = 1;
hero.Collider = new collider.Box(-12, -12, 24, 24);
hero.Transform.Position = new vector.Vector3(width / 2, height / 2);
hero.AddScript(player);

for (let i = 0; i < 5; i++) {
  const c = g.Add(game.NewObject("coin"))!;
  c.Sprite = coinLook;
  c.Collider = new collider.Circle(0, 0, 7);
  c.AddScript({ Update: (o: game.Object, dt: number) => o.Transform.Rotate(new vector.Vector3(0, 0, 1), 2 * dt) });
  place(c);
}

const runErr = g.Run(goString("sge collect (TypeScript)"), width, height);
if (runErr !== null) console.error(runErr.Error());
console.log("score:", player.score);
