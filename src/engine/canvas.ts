export type Image = EmulatedCanvas2D;

export type CanvasImageSource = {
  width: number;
  height: number;
  data: Uint8Array;
};

export type EmulatedCanvas2D = {
  width: number;
  height: number;
  getContext(contextId: "2d"): CanvasContext2D;
  toBuffer(): Uint8Array;
};

export type ImageData = {
  width: number;
  height: number;
  data: Uint8Array;
};

export type LoadImageOptions = {
  imgBytes?: Uint8Array;
  url?: string;
};

export abstract class CanvasError implements Error {
  name: string;
  message: string;
  stack?: string | undefined;
  cause?: unknown

  constructor(name: string, message: string, cause?: unknown) {
    this.name = name;
    this.message = message;
    this.name = name;
    this.message = message;
    this.cause = cause;
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ErrInvalidImageOptions);
    }
  }
}

export class ErrInvalidImageOptions extends CanvasError {
  constructor(cause?: unknown) {
    super("ErrInvalidImageOptions", "Invalid image options", cause);
  }
}

export class ErrCannotCreateFile extends CanvasError {
  constructor(cause?: unknown) {
    super("ErrCannotCreateFile", "Cannot create file", cause);
  }
}

export class ErrCannotDeleteFile extends CanvasError {
  constructor(cause?: unknown) {
    super("ErrCannotCreateFile", "Cannot create file", cause);
  }
}

export interface CanvasContext2D  {
  drawImage: (image: CanvasImageSource | EmulatedCanvas2D, x: number, y: number, width: number, height: number) => void;
  getImageData: (x: number, y: number, width: number, height: number) => ImageData;
  setFilePath(filePath: string): CanvasContext2D;
  putImageData: (imageData: ImageData, x: number, y: number) => void;
  beginPath: () => void;
  createFile(): Error | null;
  deleteFile(): Error | null;
  arc: (x: number, y: number, radius: number, startAngle: number, endAngle: number) => void;
  fill: () => void;
  fillRect: (x: number, y: number, width: number, height: number) => void;
  fillStyle: string;
};

class Context2D implements CanvasContext2D {
  private buffer: Uint8Array;
  private width: number;
  private filePath: string;
  private initialFilePath: string;
  private filePathDefined: boolean = false;
  private height: number;
  fillStyle: string = "#000000";

  constructor(width: number, height: number) {
    this.width = width;
    this.height = height;
    this.buffer = new Uint8Array(width * height * 4);
    this.filePath = "";
    this.initialFilePath = "";
  }

  setFilePath(filePath: string): CanvasContext2D {
    this.filePath = filePath;
    if (!this.filePathDefined) {
      this.initialFilePath = filePath;
      this.filePathDefined = true;
    }
    return this;
  }

  beginPath(): void {
      this.filePath = this.initialFilePath;
  }

  createFile(): Error | null {
    if (this.filePath === "") {
      return new ErrCannotCreateFile("File path not defined");
    }
    Deno.writeFileSync(this.filePath, this.buffer);
    return null;
  }

  deleteFile(): Error | null {
    if (this.filePath === "") {
      return new ErrCannotDeleteFile("File path not defined");
    }
    Deno.removeSync(this.filePath);
    return null;
  }

  arc(x: number, y: number, radius: number): void {
    for (let i = 0; i < this.height; i++) {
      for (let j = 0; j < this.width; j++) {
        const dx = j - x;
        const dy = i - y;
        if (dx * dx + dy * dy <= radius * radius) {
          const index = (i * this.width + j) * 4;
          const color = this.parseColor(this.fillStyle);
          this.buffer[index] = color.r;
          this.buffer[index + 1] = color.g;
          this.buffer[index + 2] = color.b;
          this.buffer[index + 3] = 255;
        }
      }
    }
  }

  fill(): void {
    return this.arc(0, 0, Math.min(this.width, this.height) / 2);
  }

  fillRect(x: number, y: number, width: number, height: number): void {
    const color = this.parseColor(this.fillStyle);
    for (let i = y; i < y + height; i++) {
      for (let j = x; j < x + width; j++) {
        if (i < 0 || i >= this.height || j < 0 || j >= this.width) continue;
        const index = (i * this.width + j) * 4;
        this.buffer[index] = color.r;
        this.buffer[index + 1] = color.g;
        this.buffer[index + 2] = color.b;
        this.buffer[index + 3] = 255;
      }
    }
  }

  drawImage(image: CanvasImageSource | EmulatedCanvas2D, x: number, y: number, width: number, height: number): void {
    const sourceData = image instanceof Object && 'getContext' in image
      ? image.getContext('2d').getImageData(0, 0, image.width, image.height).data
      : image.data;

    const sourceWidth = 'width' in image ? image.width : width;
    
    for (let row = 0; row < height; row++) {
      for (let col = 0; col < width; col++) {
        const sourceIndex = (Math.floor(row * image.height / height) * sourceWidth + 
                           Math.floor(col * image.width / width)) * 4;
        const targetIndex = ((y + row) * this.width + (x + col)) * 4;
        
        for (let i = 0; i < 4; i++) {
          this.buffer[targetIndex + i] = sourceData[sourceIndex + i];
        }
      }
    }
  }

  getImageData(x: number, y: number, width: number, height: number): ImageData {
    const data = new Uint8Array(width * height * 4);
    
    for (let row = 0; row < height; row++) {
      for (let col = 0; col < width; col++) {
        const sourceIndex = ((y + row) * this.width + (x + col)) * 4;
        const targetIndex = (row * width + col) * 4;
        
        for (let i = 0; i < 4; i++) {
          data[targetIndex + i] = this.buffer[sourceIndex + i];
        }
      }
    }

    return { width, height, data };
  }

  putImageData(imageData: ImageData, x: number, y: number): void {
    for (let row = 0; row < imageData.height; row++) {
      for (let col = 0; col < imageData.width; col++) {
        const sourceIndex = (row * imageData.width + col) * 4;
        const targetIndex = ((y + row) * this.width + (x + col)) * 4;
        
        for (let i = 0; i < 4; i++) {
          this.buffer[targetIndex + i] = imageData.data[sourceIndex + i];
        }
      }
    }
  }

  private parseColor(color: string): { r: number; g: number; b: number } {
    if (color.startsWith('#')) {
      return {
        r: parseInt(color.slice(1, 3), 16),
        g: parseInt(color.slice(3, 5), 16),
        b: parseInt(color.slice(5, 7), 16)
      };
    }
    // Add support for named colors
    const namedColors: Record<string, { r: number; g: number; b: number }> = {
      red: { r: 255, g: 0, b: 0 },
      blue: { r: 0, g: 0, b: 255 },
      // Add more colors as needed
    };
    return namedColors[color] || { r: 0, g: 0, b: 0 };
  }
}

export class Canvas {
  static create(width: number, height: number): EmulatedCanvas2D {
    const context = new Context2D(width, height);
    
    return {
      width,
      height,
      getContext: (_contextId: "2d") => context,
      toBuffer: () => context.getImageData(0, 0, width, height).data,
    };
  }

  static async loadImage(options: LoadImageOptions): Promise<EmulatedCanvas2D> {
    if (!options.url && !options.imgBytes) {
      throw new ErrInvalidImageOptions();
    }

    if (options.url) {
      const response = await fetch(options.url);
      const arrayBuffer = await response.arrayBuffer();
      options.imgBytes = new Uint8Array(arrayBuffer);
    }

    if (!options.imgBytes) {
      throw new ErrInvalidImageOptions();
    }

    // Create a new canvas with the image data
    const canvas = Canvas.create(1, 1); // Default size, should be adjusted based on actual image
    const ctx = canvas.getContext("2d");
    const imageData = {
      width: 1,
      height: 1,
      data: options.imgBytes
    };
    ctx.putImageData(imageData, 0, 0);
    return canvas;
  }
}

export function convertImageToUint8Array(source: CanvasImageSource | EmulatedCanvas2D): Uint8Array {
  if ('toBuffer' in source) {
    return source.toBuffer();
  }
  
  const width = source.width;
  const height = source.height;
  const canvas = Canvas.create(width, height);
  const ctx = canvas.getContext("2d");
  
  ctx.drawImage(source, 0, 0, width, height);
  

  return canvas.toBuffer();
}