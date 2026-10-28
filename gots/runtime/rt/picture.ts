// A picture the Go backend hands the window bridge, an icon or a cursor:
// a canvas's 0xAARRGGBB words, stride of them a row.

export type Picture = { pixels: Uint32Array; stride: number; width: number; height: number };

// fits reports whether p's pixels hold every row its shape says, all SDL
// reads of them.
export function fits(p: Picture): boolean {
  return p.width >= 0 && p.height >= 0 && p.stride >= p.width && p.pixels.length >= p.stride * p.height;
}

// dataURL is p as a PNG data: URL, alpha kept, for a page's icon or cursor.
export function dataURL(p: Picture): string {
  const c = (globalThis as any).document.createElement("canvas");
  c.width = p.width;
  c.height = p.height;
  const ctx = c.getContext("2d");
  const img = ctx.createImageData(p.width, p.height);
  const dst = new Uint32Array(img.data.buffer);
  let o = 0;
  for (let y = 0; y < p.height; y++) {
    for (let i = y * p.stride, end = i + p.width; i < end; i++) {
      const v = p.pixels[i];
      dst[o++] = (v & 0xff00ff00) | ((v >>> 16) & 0xff) | ((v & 0xff) << 16);
    }
  }
  ctx.putImageData(img, 0, 0);
  return c.toDataURL("image/png");
}
