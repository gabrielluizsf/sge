// The sockets behind Go's net: a read that runs out of time loses
// nothing, a woken read looks again, and the end of a stream is told.
import { connect, listen, Listener } from "../../runtime/rt/net.ts";

function assertEquals(got: unknown, want: unknown) {
  if (got !== want) throw new Error(`got ${String(got)}, want ${String(want)}`);
}

const text = (b: Uint8Array | null | "wait") => (b instanceof Uint8Array ? new TextDecoder().decode(b) : b);

Deno.test("a read past its deadline keeps what arrives for the next", async () => {
  const l = listen("127.0.0.1", 0);
  if (!(l instanceof Listener)) throw new Error(l.message);
  const [, port] = l.local();
  const [c, s] = await Promise.all([connect("127.0.0.1", port), l.accept(-1)]);
  if (s === "wait") throw new Error("accept waited");
  assertEquals(await c.read(10, 30), "wait");
  const w = s.startWrite(new TextEncoder().encode("hello world"));
  await w.wait(-1);
  assertEquals(w.done && w.n === 11, true);
  assertEquals(text(await c.read(5, -1)), "hello");
  assertEquals(text(await c.read(100, -1)), " world");
  // A reader woken to look at a new deadline returns at once.
  const r = c.read(10, -1);
  setTimeout(() => c.wake(), 10);
  assertEquals(await r, "wait");
  s.closeWrite();
  assertEquals(await c.read(10, -1), null);
  c.close();
  s.close();
  l.close();
});

Deno.test("dialing a closed port fails with ConnectionRefused", async () => {
  const l = listen("127.0.0.1", 0);
  if (!(l instanceof Listener)) throw new Error(l.message);
  const [, port] = l.local();
  l.close();
  let name = "";
  try {
    await connect("127.0.0.1", port);
  } catch (e) {
    name = (e as { name: string }).name;
  }
  assertEquals(name, "ConnectionRefused");
});

Deno.test("listening on a port in use gives AddrInUse, without throwing", () => {
  const l = listen("127.0.0.1", 0);
  if (!(l instanceof Listener)) throw new Error(l.message);
  const again = listen("127.0.0.1", l.local()[1]);
  assertEquals(again instanceof Listener ? "" : again.name, "AddrInUse");
  l.close();
});
