import {
  type CanvasImageSource,
  createCanvas,
  type EmulatedCanvas2D,
  loadImage,
} from "@josefabio/deno-canvas";

export type Image = CanvasImageSource;

export type LoadImageOptions = {
  imgBytes?: Uint8Array;
  url?: string;
};

export class ErrInvalidImageOptions implements Error {
  name: string;
  message: string;
  stack?: string | undefined;
  cause?: unknown;

  constructor(cause?: unknown) {
    this.name = "ErrInvalidImageOptions";
    this.message = "Invalid image options";
    this.cause = cause;
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ErrInvalidImageOptions);
    }
  }
}

export class Canvas {
  static create(width: number, height: number): EmulatedCanvas2D {
    return createCanvas(width, height);
  }
  static loadImage(options: LoadImageOptions): Promise<Image> {
    if (options.url) {
      return loadImage(options.url);
    }
    if (options.imgBytes) {
      return loadImage(options.imgBytes);
    }
    throw new ErrInvalidImageOptions();
  }
}

type Source = {
  width: number;
  height: number;
};

export function convertImageToUint8Array(
  source: CanvasImageSource,
): Uint8Array {
  const width = (source as Source).width || 100;
  const height = (source as Source).height || 100;

  const canvas = createCanvas(width, height);
  const ctx = canvas.getContext("2d");

  ctx.drawImage(source, 0, 0, width, height);

  const imageData = ctx.getImageData(0, 0, width, height);

  return new Uint8Array(imageData.data.buffer);
}
