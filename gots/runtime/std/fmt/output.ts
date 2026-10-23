// write sends a Go string's bytes to standard output or standard error.
export function write(s: string, toStderr = false): void {
  const g: any = globalThis;
  const bytes = new Uint8Array(s.length);
  for (let i = 0; i < s.length; i++) bytes[i] = s.charCodeAt(i);
  if (g.Deno !== undefined) {
    const out = toStderr ? g.Deno.stderr : g.Deno.stdout;
    for (let n = 0; n < bytes.length;) n += out.writeSync(bytes.subarray(n));
  } else if (g.process?.stdout !== undefined) {
    (toStderr ? g.process.stderr : g.process.stdout).write(bytes);
  } else {
    console.log(new TextDecoder().decode(bytes));
  }
}
