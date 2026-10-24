// Text crosses between TypeScript and the transpiled engine with no
// conversions written by hand. Run by TestTypeScriptBoundary, from the
// directory the engine was transpiled to.
import { canvas, game, sprite } from "./mod.ts";

function assertEquals(got: unknown, want: unknown, what: string): void {
  if (got !== want) throw new Error(`${what}: got ${JSON.stringify(got)}, want ${JSON.stringify(want)}`);
}

Deno.test("a field set from TypeScript reads back as the same text", () => {
  const o = game.NewObject("jogador");
  o.Name = "João, o rápido 🚀";
  assertEquals(o.Name, "João, o rápido 🚀", "Name");
});

Deno.test("text passed to a function is found by Go code", () => {
  const g = game.New();
  g.Add(game.NewObject("Conceição"));
  const found = g.Find("Conceição");
  assertEquals(found?.Name, "Conceição", "Find");
  assertEquals(g.Find("Conceicao"), null, "Find of other text");
});

Deno.test("an animation named in TypeScript plays by that name", async () => {
  const { animation } = await import("./mod.ts");
  const a = animation.NewAnimator();
  const [s] = sprite.Rect(1, 1, canvas.RGB(1, 2, 3));
  a.Add("pulo à frente", new animation.Animation(null, 10, true));
  assertEquals(a.Play("pulo à frente"), true, "Play");
  assertEquals(a.Current(), "pulo à frente", "Current");
  assertEquals(s !== null, true, "sprite");
});

Deno.test("text drawn from TypeScript is the same as Go's bytes", () => {
  const [cv] = canvas.NewCanvas(120, 20);
  const [cv2] = canvas.NewCanvas(120, 20);
  cv.Text(new canvas.Point(0, 0), "ação é", canvas.RGB(255, 255, 255));
  // The same text written as the UTF-8 bytes Go holds.
  const bytes = new TextEncoder().encode("ação é");
  let goBytes = "";
  for (const b of bytes) goBytes += String.fromCharCode(b < 0x80 ? b : 0xf700 + b);
  cv2.Text(new canvas.Point(0, 0), goBytes, canvas.RGB(255, 255, 255));
  let same = true;
  for (let i = 0; i < cv.Pixels!.n; i++) if (cv.Pixels!.get(i) !== cv2.Pixels!.get(i)) same = false;
  assertEquals(same, true, "pixels");
});

Deno.test("a TypeScript script's text reaches Go through an interface", () => {
  // Painter is called by Go code; a stringer made in TypeScript returns
  // JavaScript text, which Go has to read as its own.
  const g = game.New();
  const o = g.Add(game.NewObject("x"));
  let painted = "";
  o.AddScript({
    Update() {},
    Paint(obj: game.Object) {
      painted = obj.Name;
    },
  });
  o.Name = "Ícaro";
  const [cv] = canvas.NewCanvas(4, 4);
  g.Draw(cv);
  assertEquals(painted, "Ícaro", "painted name");
});

// A Go function that may block is a promise to TypeScript, the class
// method calling it included: Run paces its frames with time.Sleep.
type RunResult = ReturnType<game.Game["Run"]>;
const runIsAPromise: RunResult extends Promise<unknown> ? true : false = true;
type StepResult = ReturnType<game.Game["Step"]>;
const stepIsNot: StepResult extends Promise<unknown> ? false : true = true;
Deno.test("blocking functions are typed as promises", () => {
  if (!runIsAPromise || !stepIsNot) throw new Error("unreachable");
});
