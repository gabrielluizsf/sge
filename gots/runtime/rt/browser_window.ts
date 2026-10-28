// The window bridge in a browser: a window is a <canvas> on the page, and
// its events the page's, told to the Go backend in the shape SDL's are
// (SDL keycodes, buttons 1 to 3), so the backend reads both alike.
//
// The page may hold a <canvas data-gots> for the game to draw in; if it
// does not, one is added to the body.

import { dataURL, type Picture } from "./picture.ts";

type BridgeEvent = Record<string, unknown> & { type: string };

// SDL's keycodes for KeyboardEvent.code, but for letters and digits,
// whose keycode is their character.
const keycodes: Record<string, number> = {
  Enter: 13, NumpadEnter: 13, Escape: 27, Backspace: 8, Tab: 9, Delete: 127, Space: 32,
  Quote: 39, Comma: 44, Minus: 45, Period: 46, Slash: 47, Semicolon: 59, Equal: 61,
  BracketLeft: 91, Backslash: 92, BracketRight: 93, Backquote: 96,
  CapsLock: 0x40000039, Insert: 0x40000049, Home: 0x4000004a, PageUp: 0x4000004b,
  End: 0x4000004d, PageDown: 0x4000004e, ArrowRight: 0x4000004f, ArrowLeft: 0x40000050,
  ArrowDown: 0x40000051, ArrowUp: 0x40000052,
  ControlLeft: 0x400000e0, ShiftLeft: 0x400000e1, AltLeft: 0x400000e2, MetaLeft: 0x400000e3,
  ControlRight: 0x400000e4, ShiftRight: 0x400000e5, AltRight: 0x400000e6, MetaRight: 0x400000e7,
};

// sdlKey is the SDL keycode of the key code names, or 0.
export function sdlKey(code: string): number {
  if (/^Key[A-Z]$/.test(code)) return code.charCodeAt(3) + 32;
  if (/^Digit[0-9]$/.test(code)) return code.charCodeAt(5);
  const f = /^F([1-9]|1[0-2])$/.exec(code);
  if (f) return 0x4000003a + Number(f[1]) - 1;
  return keycodes[code] ?? 0;
}

// sdlMods are SDL's modifier bits for what e holds down.
export function sdlMods(e: { shiftKey: boolean; ctrlKey: boolean; altKey: boolean; metaKey: boolean }): number {
  return (e.shiftKey ? 0x0001 : 0) | (e.ctrlKey ? 0x0040 : 0) | (e.altKey ? 0x0100 : 0) | (e.metaKey ? 0x0400 : 0);
}

// swizzle copies the canvas's 0xAARRGGBB words, stride a row, to an
// ImageData's bytes, R G B A, which a little-endian Uint32Array reads as
// 0xAABBGGRR: red and blue trade places, and the window is opaque, as
// SDL's texture is drawn without blending.
export function swizzle(src: Uint32Array, stride: number, width: number, height: number, dst: Uint32Array): void {
  let o = 0;
  for (let y = 0; y < height; y++) {
    let i = y * stride;
    for (const end = i + width; i < end; i++) {
      const p = src[i];
      dst[o++] = 0xff000000 | (p & 0x0000ff00) | ((p >>> 16) & 0xff) | ((p & 0xff) << 16);
    }
  }
}

export class BrowserWindow {
  private events: BridgeEvent[] = [];
  private canvas: any;
  private frame: any; // what goes full screen: the canvas centred on black
  private ctx: any;
  private image: any = null;
  private pasted = "";
  private off: (() => void)[] = [];

  constructor(title: string, width: number, height: number) {
    const doc = (globalThis as any).document;
    doc.title = title;
    this.canvas = doc.querySelector("canvas[data-gots]");
    if (this.canvas === null) {
      this.canvas = doc.createElement("canvas");
      this.canvas.setAttribute("data-gots", "");
      doc.body.appendChild(this.canvas);
    }
    this.frame = this.canvas.parentElement;
    this.canvas.width = width;
    this.canvas.height = height;
    // Pixel art stays sharp when the page scales the canvas.
    this.canvas.style.imageRendering = "pixelated";
    this.canvas.tabIndex = 0;
    this.ctx = this.canvas.getContext("2d", { alpha: false });
    this.listen();
    this.canvas.focus();
  }

  private on(target: any, type: string, f: (e: any) => void, opts?: unknown): void {
    target.addEventListener(type, f, opts);
    this.off.push(() => target.removeEventListener(type, f, opts));
  }

  private listen(): void {
    const g = globalThis as any, c = this.canvas;
    const key = (type: string) => (e: any) => {
      const sym = sdlKey(e.code);
      if (sym !== 0) {
        this.events.push({ type, sym, mods: sdlMods(e), repeat: e.repeat === true });
        // The game's keys are the game's, not the page's to scroll with;
        // the browser's shortcuts stay the browser's.
        if (!e.ctrlKey && !e.metaKey) e.preventDefault();
      }
      if (type === "keydown" && typeof e.key === "string" && [...e.key].length === 1 && !e.ctrlKey && !e.metaKey) {
        this.events.push({ type: "text", text: e.key });
      }
    };
    this.on(c, "keydown", key("keydown"));
    this.on(c, "keyup", key("keyup"));
    this.on(c, "mousemove", (e: any) => {
      const [x, y] = this.point(e);
      this.events.push({ type: "move", x, y, dx: Math.round(e.movementX ?? 0), dy: Math.round(e.movementY ?? 0) });
    });
    const button = (type: string) => (e: any) => {
      if (e.button > 2) return;
      if (type === "down") c.focus();
      const [x, y] = this.point(e);
      // DOM's buttons count from 0 as left, middle, right; SDL's from 1.
      this.events.push({ type, button: e.button + 1, x, y });
    };
    this.on(c, "mousedown", button("down"));
    this.on(c, "mouseup", button("up"));
    this.on(c, "contextmenu", (e: any) => e.preventDefault());
    this.on(c, "wheel", (e: any) => {
      // SDL's wheel is positive away from the user, DOM's deltaY toward.
      if (e.deltaY !== 0) this.events.push({ type: "wheel", wheel: e.deltaY < 0 ? 1 : -1 });
      e.preventDefault();
    }, { passive: false });
    this.on(c, "focus", () => this.events.push({ type: "focus", focused: true }));
    this.on(c, "blur", () => this.events.push({ type: "focus", focused: false }));
    this.on(c, "paste", (e: any) => {
      this.pasted = e.clipboardData?.getData("text") ?? this.pasted;
    });
    this.on(g, "pagehide", () => this.events.push({ type: "close" }));
    this.on(g.document, "visibilitychange", () => {
      if (g.document.visibilityState === "visible") this.events.push({ type: "expose" });
    });
    this.on(g.document, "fullscreenchange", () => this.fit());
    this.on(g, "resize", () => this.fit());
  }

  // point is where e is on the canvas, in its pixels, however the page
  // scales it.
  private point(e: any): [number, number] {
    const r = this.canvas.getBoundingClientRect();
    const sx = r.width > 0 ? this.canvas.width / r.width : 1;
    const sy = r.height > 0 ? this.canvas.height / r.height : 1;
    return [Math.floor((e.clientX - r.left) * sx), Math.floor((e.clientY - r.top) * sy)];
  }

  // fit scales the canvas to fill the screen, keeping its shape, while
  // full screen, and back to its own size after.
  private fit(): void {
    const g = globalThis as any, s = this.canvas.style;
    if (g.document.fullscreenElement !== this.frame) {
      s.width = s.height = "";
      return;
    }
    const k = Math.min(g.innerWidth / this.canvas.width, g.innerHeight / this.canvas.height);
    s.width = Math.floor(this.canvas.width * k) + "px";
    s.height = Math.floor(this.canvas.height * k) + "px";
  }

  present(pixels: Uint32Array, stride: number, width: number, height: number): boolean {
    if (this.canvas.width !== width || this.canvas.height !== height) {
      this.canvas.width = width;
      this.canvas.height = height;
      this.fit();
    }
    if (this.image === null || this.image.width !== width || this.image.height !== height) {
      this.image = this.ctx.createImageData(width, height);
    }
    swizzle(pixels, stride, width, height, new Uint32Array(this.image.data.buffer));
    this.ctx.putImageData(this.image, 0, 0);
    return true;
  }

  poll(): BridgeEvent[] {
    const out = this.events;
    this.events = [];
    return out;
  }

  setTitle(title: string): void {
    (globalThis as any).document.title = title;
  }

  // setFullscreen asks the page for full screen, which a browser grants
  // only in answer to the player's click or key.
  setFullscreen(on: boolean): boolean {
    const doc = (globalThis as any).document;
    if (on) {
      if (typeof this.frame?.requestFullscreen !== "function") return false;
      this.frame.style.background = "black";
      this.frame.requestFullscreen().catch(() => {});
    } else if (doc.fullscreenElement) {
      doc.exitFullscreen().catch(() => {});
    }
    return true;
  }

  setOpacity(alpha: number): boolean {
    this.canvas.style.opacity = String(alpha / 255);
    return true;
  }

  setSize(width: number, height: number): boolean {
    this.canvas.width = width;
    this.canvas.height = height;
    this.events.push({ type: "resize", width, height });
    return true;
  }

  setLimits(): void {}

  // setIcon makes the picture the page's icon.
  setIcon(p: Picture): boolean {
    const doc = (globalThis as any).document;
    let link = doc.querySelector("link[rel~='icon']");
    if (link === null) {
      link = doc.createElement("link");
      link.rel = "icon";
      doc.head.appendChild(link);
    }
    link.href = dataURL(p);
    return true;
  }

  // setCursor shows p as the mouse's cursor over the canvas, its hot spot
  // at x, y, or the page's own for none.
  setCursor(p: Picture | null, x: number, y: number): boolean {
    this.canvas.style.cursor = p === null ? "" : `url(${dataURL(p)}) ${x} ${y}, auto`;
    return true;
  }

  scale(): number {
    return (globalThis as any).devicePixelRatio ?? 0;
  }

  // notify shows a notification, asking the player first if the page may,
  // as a browser asks.
  notify(title: string, body: string): boolean {
    const N = (globalThis as any).Notification;
    if (typeof N !== "function" || N.permission === "denied") return false;
    const show = () => new N(title, { body });
    if (N.permission === "granted") show();
    else N.requestPermission().then((p: string) => p === "granted" && show(), () => {});
    return true;
  }

  systemDark(): [boolean, boolean] {
    const m = (globalThis as any).matchMedia;
    return typeof m === "function" ? [m("(prefers-color-scheme: dark)").matches, true] : [false, false];
  }

  display(): number[] | null {
    const s = (globalThis as any).screen;
    return s ? [s.width, s.height, 60] : null;
  }

  // clipboard is what was last pasted on the canvas: a page reads the
  // clipboard only when the player pastes.
  clipboard(): string {
    return this.pasted;
  }

  setClipboard(text: string): boolean {
    const c = (globalThis as any).navigator?.clipboard;
    if (typeof c?.writeText !== "function") return false;
    c.writeText(text).catch(() => {});
    return true;
  }

  close(): void {
    for (const f of this.off) f();
    this.off = [];
  }
}
