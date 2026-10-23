// fetchBytes reads an address whole with fetch, holding the limits a Go
// program's own read would: an answer other than 200 is refused, so is a
// body that says or turns out to be longer than max bytes, read no further
// than that, and so is a read not finished within timeoutMs. It is put on
// globalThis as gotsFetchBytes for overlays reaching it through syscall/js.
export async function fetchBytes(url: string, max: number, timeoutMs: number): Promise<Uint8Array> {
  const resp = await fetch(url, { signal: AbortSignal.timeout(timeoutMs) });
  const tooBig = () => new Error(`it is bigger than ${max} bytes`);
  if (resp.status !== 200) {
    await resp.body?.cancel();
    throw new Error(`it answered ${resp.status} ${resp.statusText}`);
  }
  if (Number(resp.headers.get("content-length")) > max) {
    await resp.body?.cancel();
    throw tooBig();
  }
  if (resp.body === null) return new Uint8Array(0);
  const reader = resp.body.getReader();
  const chunks: Uint8Array[] = [];
  let n = 0;
  for (;;) {
    const { done, value } = await reader.read();
    if (done) break;
    n += value.length;
    if (n > max) {
      await reader.cancel();
      throw tooBig();
    }
    chunks.push(value);
  }
  const out = new Uint8Array(n);
  let at = 0;
  for (const c of chunks) {
    out.set(c, at);
    at += c.length;
  }
  return out;
}

(globalThis as any).gotsFetchBytes = fetchBytes;
