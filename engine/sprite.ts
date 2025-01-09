import { Image, Canvas } from "./canvas.ts";

// Sprite class represents a 2D sprite with position, size, and image
export class Sprite {
  x: number = 0;
  y: number = 0;
  width: number;
  height: number;
  image: Image;

  constructor() {
    this.width = 100;
    this.height = 100;
    this.image = this.createSquare(100, "blue");
  }

  // Set the position of the sprite
  setPosition(x: number, y: number): void {
    this.x = x;
    this.y = y;
  }

  // Set the size of the sprite
  setSize(width: number, height: number): void {
    this.width = width;
    this.height = height;
  }

  // Set the image of the sprite from an image file or buffer
  async setImage(imgBytes: Uint8Array): Promise<void> {
    this.image = await Canvas.loadImage({ imgBytes });
  }

  // Change the color of the sprite's image
  changeColor(newColor: string): void {
    const canvas = Canvas.create(this.width, this.height);
    const ctx = canvas.getContext("2d")!;
    ctx.drawImage(this.image, 0, 0, this.width, this.height);

    // Create a new image with the desired color
    const imageData = ctx.getImageData(0, 0, this.width, this.height);
    const data = imageData.data;

    for (let i = 0; i < data.length; i += 4) {
      const alpha = data[i + 3];
      if (alpha !== 0) {
        // Change the RGB values based on the new color
        data[i] = parseInt(newColor.substring(1, 3), 16); // Red
        data[i + 1] = parseInt(newColor.substring(3, 5), 16); // Green
        data[i + 2] = parseInt(newColor.substring(5, 7), 16); // Blue
      }
    }

    ctx.putImageData(imageData, 0, 0);
    this.image = canvas;
  }

  // Move the sprite by a certain amount
  move(dx: number, dy: number): void {
    this.x += dx;
    this.y += dy;
  }

  // Resize the sprite
  resize(newWidth: number, newHeight: number): void {
    this.width = newWidth;
    this.height = newHeight;
  }

  // Create a circle image
  createCircle(radius: number, color: string): Image {
    const canvas = Canvas.create(radius * 2, radius * 2);
    const ctx = canvas.getContext("2d")!;
    ctx.beginPath();
    ctx.arc(radius, radius, radius, 0, Math.PI * 2);
    ctx.fillStyle = color;
    ctx.fill();
    return canvas;
  }

  // Create a square image
  createSquare(side: number, color: string): Image {
    const canvas = Canvas.create(side, side);
    const ctx = canvas.getContext("2d")!;
    ctx.fillStyle = color;
    ctx.fillRect(0, 0, side, side);
    return canvas;
  }
}
