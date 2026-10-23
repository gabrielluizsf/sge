// fetchBytes is how svg's overlay reads a picture from the network; it has
// to hold the limits antui's own net/http read holds: only a 200, nothing
// past the size limit in memory, and nothing slower than the timeout.
import { fetchBytes } from "../../runtime/rt/fetch.ts";

function serve(handler: (req: Request) => Response): { url: string; close: () => Promise<void> } {
  const server = Deno.serve({ port: 0, hostname: "127.0.0.1", onListen() {} }, handler);
  return { url: `http://127.0.0.1:${server.addr.port}/`, close: () => server.shutdown() };
}

async function refusal(p: Promise<unknown>): Promise<string> {
  try {
    await p;
  } catch (e) {
    return e instanceof Error ? e.message : String(e);
  }
  throw new Error("it was read");
}

Deno.test("a picture within the limit is read whole", async () => {
  const s = serve(() => new Response(new Uint8Array([1, 2, 3])));
  try {
    const got = await fetchBytes(s.url, 3, 5000);
    if (got.join() !== "1,2,3") throw new Error(`got ${got}`);
  } finally {
    await s.close();
  }
});

Deno.test("an answer other than 200 is refused", async () => {
  const s = serve(() => new Response("gone", { status: 404, statusText: "Not Found" }));
  try {
    const msg = await refusal(fetchBytes(s.url, 1024, 5000));
    if (!msg.includes("404")) throw new Error(msg);
  } finally {
    await s.close();
  }
});

Deno.test("a declared length past the limit is refused before reading", async () => {
  const s = serve(() => new Response("x".repeat(2048)));
  try {
    const msg = await refusal(fetchBytes(s.url, 1024, 5000));
    if (!msg.includes("bigger than 1024")) throw new Error(msg);
  } finally {
    await s.close();
  }
});

Deno.test("an endless body is cut off at the limit", async () => {
  let sent = 0;
  const s = serve(() =>
    new Response(
      new ReadableStream({
        pull(c) {
          sent += 4096;
          c.enqueue(new Uint8Array(4096));
        },
      }),
    )
  );
  try {
    const msg = await refusal(fetchBytes(s.url, 10000, 5000));
    if (!msg.includes("bigger than 10000")) throw new Error(msg);
    // The server runs ahead into socket buffers; what matters is that the
    // read stopped, where reading the whole body never would.
    if (sent > 16 << 20) throw new Error(`the server sent ${sent} bytes past a 10000 byte limit`);
  } finally {
    await s.close();
  }
});

Deno.test("a server that never finishes is given up on", async () => {
  let stop = () => {};
  const s = serve(() =>
    new Response(
      new ReadableStream({
        start(c) {
          c.enqueue(new Uint8Array(1));
          stop = () => c.close();
        },
      }),
    )
  );
  try {
    const msg = await refusal(fetchBytes(s.url, 1024, 200));
    if (msg === "it was read") throw new Error(msg);
  } finally {
    stop();
    await s.close();
  }
});
