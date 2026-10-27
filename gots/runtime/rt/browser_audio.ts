// A stream played with Web Audio: each block a buffer started right
// where the one before ends. A page may not make a sound before the one
// looking at it does something, so the context starts when they first
// press a key or the pointer; until then nothing plays and the queue,
// whose time stands still, stays as full as the Go side keeps it.

export class BrowserAudio {
  ctx: any;
  rate: number;
  next = 0;

  constructor(rate: number) {
    const G: any = globalThis;
    this.ctx = new G.AudioContext({ sampleRate: rate });
    this.rate = this.ctx.sampleRate;
    const start = () => this.ctx.resume();
    for (const ev of ["pointerdown", "keydown"]) G.addEventListener?.(ev, start, { once: true });
  }

  queued(): number {
    return Math.max(0, this.next - this.ctx.currentTime);
  }

  push(block: Float32Array): void {
    const frames = block.length >> 1;
    if (frames === 0) return;
    const buf = this.ctx.createBuffer(2, frames, this.rate);
    const l = buf.getChannelData(0), r = buf.getChannelData(1);
    for (let i = 0; i < frames; i++) {
      l[i] = block[2 * i];
      r[i] = block[2 * i + 1];
    }
    const src = this.ctx.createBufferSource();
    src.buffer = buf;
    src.connect(this.ctx.destination);
    this.next = Math.max(this.next, this.ctx.currentTime);
    src.start(this.next);
    this.next += frames / this.rate;
  }

  close(): void {
    this.ctx.close();
  }
}
