// The window bridge: what a transpiled program's window backend calls to
// get a native window on Deno, through SDL2 and Deno's FFI. It is put on
// globalThis as gotsWindow, and SDL2 is only loaded when a window opens,
// so programs that draw off screen need neither SDL2 nor --allow-ffi.

const SDL_INIT_VIDEO = 0x20;
const SDL_WINDOWPOS_CENTERED = 0x2fff0000;
const SDL_WINDOW_SHOWN = 0x4, SDL_WINDOW_RESIZABLE = 0x20;
const SDL_WINDOW_FULLSCREEN_DESKTOP = 0x1001;
const SDL_RENDERER_ACCELERATED = 0x2;
const SDL_PIXELFORMAT_ARGB8888 = 0x16362004;
const SDL_TEXTUREACCESS_STREAMING = 1;

const SDL_QUIT = 0x100, SDL_WINDOWEVENT = 0x200, SDL_KEYDOWN = 0x300, SDL_KEYUP = 0x301;
const SDL_TEXTINPUT = 0x303, SDL_MOUSEMOTION = 0x400, SDL_MOUSEBUTTONDOWN = 0x401;
const SDL_MOUSEBUTTONUP = 0x402, SDL_MOUSEWHEEL = 0x403, SDL_DROPFILE = 0x1000;
const WINDOWEVENT_EXPOSED = 3, WINDOWEVENT_SIZE_CHANGED = 6, WINDOWEVENT_FOCUS_GAINED = 12;
const WINDOWEVENT_FOCUS_LOST = 13, WINDOWEVENT_CLOSE = 14;

const symbols = {
  SDL_Init: { parameters: ["u32"], result: "i32" },
  SDL_GetError: { parameters: [], result: "pointer" },
  SDL_CreateWindow: { parameters: ["buffer", "i32", "i32", "i32", "i32", "u32"], result: "pointer" },
  SDL_CreateRenderer: { parameters: ["pointer", "i32", "u32"], result: "pointer" },
  SDL_CreateTexture: { parameters: ["pointer", "u32", "i32", "i32", "i32"], result: "pointer" },
  SDL_DestroyTexture: { parameters: ["pointer"], result: "void" },
  SDL_DestroyRenderer: { parameters: ["pointer"], result: "void" },
  SDL_DestroyWindow: { parameters: ["pointer"], result: "void" },
  SDL_UpdateTexture: { parameters: ["pointer", "pointer", "buffer", "i32"], result: "i32" },
  SDL_RenderCopy: { parameters: ["pointer", "pointer", "pointer", "pointer"], result: "i32" },
  SDL_RenderPresent: { parameters: ["pointer"], result: "void" },
  SDL_PollEvent: { parameters: ["buffer"], result: "i32" },
  SDL_SetWindowTitle: { parameters: ["pointer", "buffer"], result: "void" },
  SDL_SetWindowFullscreen: { parameters: ["pointer", "u32"], result: "i32" },
  SDL_SetWindowOpacity: { parameters: ["pointer", "f32"], result: "i32" },
  SDL_SetWindowSize: { parameters: ["pointer", "i32", "i32"], result: "void" },
  SDL_SetWindowMinimumSize: { parameters: ["pointer", "i32", "i32"], result: "void" },
  SDL_SetWindowMaximumSize: { parameters: ["pointer", "i32", "i32"], result: "void" },
  SDL_GetDesktopDisplayMode: { parameters: ["i32", "buffer"], result: "i32" },
  SDL_GetClipboardText: { parameters: [], result: "pointer" },
  SDL_SetClipboardText: { parameters: ["buffer"], result: "i32" },
  SDL_free: { parameters: ["pointer"], result: "void" },
  SDL_StartTextInput: { parameters: [], result: "void" },
  SDL_Quit: { parameters: [], result: "void" },
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
  if (D?.dlopen === undefined) throw new Error("opening a window needs Deno");
  let lastErr: unknown = null;
  for (const name of libraryNames[D.build.os] ?? libraryNames.linux) {
    try {
      sdl = D.dlopen(name, symbols).symbols;
      break;
    } catch (e) {
      lastErr = e;
    }
  }
  if (sdl === null) throw new Error("SDL2 could not be loaded (" + String(lastErr) + "); install SDL2 and run Deno with --allow-ffi");
  if (sdl.SDL_Init(SDL_INIT_VIDEO) !== 0) throw new Error("SDL_Init: " + sdlError());
  return sdl;
}

function cstr(s: string): Uint8Array {
  return new TextEncoder().encode(s + "\0");
}

function readCString(p: any): string {
  if (p === null) return "";
  return new (globalThis as any).Deno.UnsafePointerView(p).getCString();
}

function sdlError(): string {
  return readCString(sdl.SDL_GetError());
}

// BridgeEvent is one window event, in the shape the Go backend reads.
type BridgeEvent = {
  type: string;
  sym?: number;
  mods?: number;
  repeat?: boolean;
  x?: number;
  y?: number;
  dx?: number;
  dy?: number;
  wheel?: number;
  button?: number;
  width?: number;
  height?: number;
  focused?: boolean;
  text?: string;
  files?: string[];
};

class NativeWindow {
  win: any;
  renderer: any;
  texture: any = null;
  texW = 0;
  texH = 0;
  event = new Uint8Array(56);
  view = new DataView(this.event.buffer);

  constructor(title: string, width: number, height: number) {
    const s = load();
    this.win = s.SDL_CreateWindow(cstr(title), SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, width, height, SDL_WINDOW_SHOWN | SDL_WINDOW_RESIZABLE);
    if (this.win === null) throw new Error("SDL_CreateWindow: " + sdlError());
    this.renderer = s.SDL_CreateRenderer(this.win, -1, SDL_RENDERER_ACCELERATED);
    if (this.renderer === null) this.renderer = s.SDL_CreateRenderer(this.win, -1, 0);
    if (this.renderer === null) throw new Error("SDL_CreateRenderer: " + sdlError());
    s.SDL_StartTextInput();
  }

  // present shows the canvas, pixels being its 0xAARRGGBB words, stride
  // of them a row.
  present(pixels: Uint32Array, stride: number, width: number, height: number): boolean {
    const s = sdl;
    if (this.texture === null || this.texW !== width || this.texH !== height) {
      if (this.texture !== null) s.SDL_DestroyTexture(this.texture);
      this.texture = s.SDL_CreateTexture(this.renderer, SDL_PIXELFORMAT_ARGB8888, SDL_TEXTUREACCESS_STREAMING, width, height);
      this.texW = width;
      this.texH = height;
      if (this.texture === null) return false;
    }
    if (s.SDL_UpdateTexture(this.texture, null, pixels, stride * 4) !== 0) return false;
    s.SDL_RenderCopy(this.renderer, this.texture, null, null);
    s.SDL_RenderPresent(this.renderer);
    return true;
  }

  poll(): BridgeEvent[] {
    const s = sdl, v = this.view;
    const out: BridgeEvent[] = [];
    while (s.SDL_PollEvent(this.event) !== 0) {
      const type = v.getUint32(0, true);
      switch (type) {
        case SDL_QUIT:
          out.push({ type: "close" });
          break;
        case SDL_WINDOWEVENT: {
          const ev = v.getUint8(12);
          if (ev === WINDOWEVENT_SIZE_CHANGED) out.push({ type: "resize", width: v.getInt32(16, true), height: v.getInt32(20, true) });
          else if (ev === WINDOWEVENT_EXPOSED) out.push({ type: "expose" });
          else if (ev === WINDOWEVENT_FOCUS_GAINED || ev === WINDOWEVENT_FOCUS_LOST) out.push({ type: "focus", focused: ev === WINDOWEVENT_FOCUS_GAINED });
          else if (ev === WINDOWEVENT_CLOSE) out.push({ type: "close" });
          break;
        }
        case SDL_KEYDOWN:
        case SDL_KEYUP:
          out.push({
            type: type === SDL_KEYDOWN ? "keydown" : "keyup",
            repeat: v.getUint8(13) !== 0,
            sym: v.getInt32(20, true),
            mods: v.getUint16(24, true),
          });
          break;
        case SDL_TEXTINPUT: {
          const bytes = this.event.subarray(12, 44);
          const end = bytes.indexOf(0);
          out.push({ type: "text", text: new TextDecoder().decode(bytes.subarray(0, end < 0 ? 32 : end)) });
          break;
        }
        case SDL_MOUSEMOTION:
          out.push({ type: "move", x: v.getInt32(20, true), y: v.getInt32(24, true), dx: v.getInt32(28, true), dy: v.getInt32(32, true) });
          break;
        case SDL_MOUSEBUTTONDOWN:
        case SDL_MOUSEBUTTONUP:
          out.push({ type: type === SDL_MOUSEBUTTONDOWN ? "down" : "up", button: v.getUint8(16), x: v.getInt32(20, true), y: v.getInt32(24, true) });
          break;
        case SDL_MOUSEWHEEL:
          out.push({ type: "wheel", wheel: v.getInt32(20, true) });
          break;
        case SDL_DROPFILE: {
          const D: any = (globalThis as any).Deno;
          const p = D.UnsafePointer.create(v.getBigUint64(8, true));
          out.push({ type: "drop", files: [readCString(p)] });
          s.SDL_free(p);
          break;
        }
      }
    }
    return out;
  }

  setTitle(title: string): void {
    sdl.SDL_SetWindowTitle(this.win, cstr(title));
  }

  setFullscreen(on: boolean): boolean {
    return sdl.SDL_SetWindowFullscreen(this.win, on ? SDL_WINDOW_FULLSCREEN_DESKTOP : 0) === 0;
  }

  setOpacity(alpha: number): boolean {
    return sdl.SDL_SetWindowOpacity(this.win, alpha / 255) === 0;
  }

  setSize(width: number, height: number): boolean {
    sdl.SDL_SetWindowSize(this.win, width, height);
    return true;
  }

  setLimits(minW: number, minH: number, maxW: number, maxH: number): void {
    sdl.SDL_SetWindowMinimumSize(this.win, minW, minH);
    if (maxW > 0 && maxH > 0) sdl.SDL_SetWindowMaximumSize(this.win, maxW, maxH);
  }

  // display is the desktop's size and refresh rate, or null.
  display(): number[] | null {
    const mode = new Uint8Array(24);
    if (sdl.SDL_GetDesktopDisplayMode(0, mode) !== 0) return null;
    const v = new DataView(mode.buffer);
    return [v.getInt32(4, true), v.getInt32(8, true), v.getInt32(12, true)];
  }

  clipboard(): string {
    const p = sdl.SDL_GetClipboardText();
    const s = readCString(p);
    if (p !== null) sdl.SDL_free(p);
    return s;
  }

  setClipboard(text: string): boolean {
    return sdl.SDL_SetClipboardText(cstr(text)) === 0;
  }

  close(): void {
    if (this.texture !== null) sdl.SDL_DestroyTexture(this.texture);
    sdl.SDL_DestroyRenderer(this.renderer);
    sdl.SDL_DestroyWindow(this.win);
    this.texture = null;
  }
}

// open is what the Go backend calls: a window, or the reason there is
// none as a string.
function open(title: string, width: number, height: number): NativeWindow | string {
  try {
    return new NativeWindow(title, width, height);
  } catch (e) {
    return e instanceof Error ? e.message : String(e);
  }
}

(globalThis as any).gotsWindow = { open };
