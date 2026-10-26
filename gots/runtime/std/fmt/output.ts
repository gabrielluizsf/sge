// write sends a Go string's bytes to standard output or standard error.
export function write(s: string, toStderr = false): void {
  const g: any = globalThis;
  const bytes = new Uint8Array(s.length);
  for (let i = 0; i < s.length; i++) {
    const c = s.charCodeAt(i);
    bytes[i] = c >= 0xf700 ? c - 0xf700 : c;
  }
  if (g.Deno !== undefined) {
    const out = toStderr ? g.Deno.stderr : g.Deno.stdout;
    for (let n = 0; n < bytes.length;) n += out.writeSync(bytes.subarray(n));
  } else if (g.process?.stdout !== undefined) {
    (toStderr ? g.process.stderr : g.process.stdout).write(bytes);
  } else {
    toConsole(bytes, toStderr);
  }
}

// A console, as a browser has, takes lines: what is written is held until
// its line ends, or for a moment of quiet, so Print without a newline
// does not make a line of each call, nor split a character in two.
const held: Uint8Array[] = [new Uint8Array(0), new Uint8Array(0)];
const timers: (ReturnType<typeof setTimeout> | undefined)[] = [undefined, undefined];

function toConsole(bytes: Uint8Array, toStderr: boolean): void {
  const k = toStderr ? 1 : 0;
  const all = new Uint8Array(held[k].length + bytes.length);
  all.set(held[k]);
  all.set(bytes, held[k].length);
  const nl = all.lastIndexOf(10);
  if (nl >= 0) {
    emit(all.subarray(0, nl), toStderr);
    held[k] = all.slice(nl + 1);
  } else {
    held[k] = all;
  }
  clearTimeout(timers[k]);
  timers[k] = undefined;
  if (held[k].length > 0) {
    timers[k] = setTimeout(() => {
      emit(held[k], toStderr);
      held[k] = new Uint8Array(0);
    }, 100);
  }
}

function emit(bytes: Uint8Array, toStderr: boolean): void {
  const text = new TextDecoder().decode(bytes);
  if (toStderr) console.error(text);
  else console.log(text);
}
