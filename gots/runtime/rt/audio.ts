// The audio bridge: what a transpiled program's audio calls to play its
// stream, through SDL2's queue and Deno's FFI, or Web Audio in a browser
// (see browser_audio.ts). It is put on globalThis as gotsAudio, and SDL2
// is only loaded when a stream opens.
//
// A stream only queues what it is given, as float32 stereo frames; the
// Go side keeps a little of it queued, and asks how much is. SDL2 is
// asked for exactly that format, and a device that would play another is
// refused rather than handed bytes it would read wrong.

import { BrowserAudio } from "./browser_audio.ts";

const SDL_INIT_AUDIO = 0x10;
const AUDIO_F32LSB = 0x8120;
const SDL_AUDIO_ALLOW_FREQUENCY_CHANGE = 0x1;
const specSize = 32; // SDL_AudioSpec on a 64-bit machine

const symbols = {
  SDL_InitSubSystem: { parameters: ["u32"], result: "i32" },
  SDL_GetError: { parameters: [], result: "pointer" },
  SDL_OpenAudioDevice: { parameters: ["pointer", "i32", "buffer", "buffer", "i32"], result: "u32" },
  SDL_QueueAudio: { parameters: ["u32", "buffer", "u32"], result: "i32" },
  SDL_GetQueuedAudioSize: { parameters: ["u32"], result: "u32" },
  SDL_PauseAudioDevice: { parameters: ["u32", "i32"], result: "void" },
  SDL_CloseAudioDevice: { parameters: ["u32"], result: "void" },
} as const;

const libraryNames: Record<string, string[]> = {
  linux: ["libSDL2-2.0.so.0", "libSDL2.so"],
  darwin: ["libSDL2-2.0.0.dylib", "libSDL2.dylib", "/opt/homebrew/lib/libSDL2.dylib", "/usr/local/lib/libSDL2.dylib"],
  windows: ["SDL2.dll"],
};

let sdl: any = null;

function load(): any {
  if (sdl !== null) return sdl;
  const D: any = (globalThis as any).Deno;
  for (const name of libraryNames[D.build.os] ?? libraryNames.linux) {
    try {
      sdl = D.dlopen(name, symbols).symbols;
      break;
    } catch {
      // the next name
    }
  }
  if (sdl === null) throw new Error("SDL2 could not be loaded; install SDL2 and run Deno with --allow-ffi");
  if (sdl.SDL_InitSubSystem(SDL_INIT_AUDIO) !== 0) throw new Error("SDL_InitSubSystem: " + sdlError());
  return sdl;
}

function sdlError(): string {
  const p = sdl.SDL_GetError();
  return p === null ? "" : new (globalThis as any).Deno.UnsafePointerView(p).getCString();
}

// spec is an SDL_AudioSpec asking for float32 stereo at rate, in blocks of
// samples frames.
export function spec(rate: number, samples: number): Uint8Array {
  const b = new Uint8Array(specSize);
  const v = new DataView(b.buffer);
  v.setInt32(0, rate, true);
  v.setUint16(4, AUDIO_F32LSB, true);
  v.setUint8(6, 2);
  v.setUint16(8, samples, true);
  return b;
}

// accepted is the rate of the device obtained says SDL opened, or 0 when it
// is not the float32 stereo asked for.
export function accepted(obtained: Uint8Array): number {
  const v = new DataView(obtained.buffer, obtained.byteOffset, obtained.byteLength);
  const rate = v.getInt32(0, true);
  if (v.getUint16(4, true) !== AUDIO_F32LSB || v.getUint8(6) !== 2 || rate < 8000 || rate > 384000) return 0;
  return rate;
}

class NativeAudio {
  dev: number;
  rate: number;
  frameBytes = 8;

  constructor(rate: number) {
    const s = load();
    const obtained = new Uint8Array(specSize);
    this.dev = s.SDL_OpenAudioDevice(null, 0, spec(rate, 512), obtained, SDL_AUDIO_ALLOW_FREQUENCY_CHANGE);
    if (this.dev === 0) throw new Error("SDL_OpenAudioDevice: " + sdlError());
    this.rate = accepted(obtained);
    if (this.rate === 0) {
      s.SDL_CloseAudioDevice(this.dev);
      throw new Error("the audio device would not play float32 stereo");
    }
    s.SDL_PauseAudioDevice(this.dev, 0);
  }

  // queued is how many seconds of sound are waiting to be played.
  queued(): number {
    return sdl.SDL_GetQueuedAudioSize(this.dev) / (this.rate * this.frameBytes);
  }

  // push queues block, interleaved stereo frames.
  push(block: Float32Array): void {
    const bytes = new Uint8Array(block.buffer, block.byteOffset, block.byteLength - block.byteLength % this.frameBytes);
    sdl.SDL_QueueAudio(this.dev, bytes, bytes.byteLength);
  }

  close(): void {
    if (this.dev !== 0) sdl.SDL_CloseAudioDevice(this.dev);
    this.dev = 0;
  }
}

// open opens a stream of float32 stereo at about rate frames a second, or
// gives back why it cannot.
function open(_name: string, rate: number): NativeAudio | BrowserAudio | string {
  try {
    if ((globalThis as any).Deno?.dlopen !== undefined) return new NativeAudio(rate);
    if ((globalThis as any).AudioContext !== undefined) return new BrowserAudio(rate);
    return "this runtime has no way to play sound";
  } catch (e) {
    return String(e instanceof Error ? e.message : e);
  }
}

(globalThis as any).gotsAudio = { open };
