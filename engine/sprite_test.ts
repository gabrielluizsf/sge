import { assertEquals, assertNotEquals } from "https://deno.land/std@0.161.0/testing/asserts.ts";
import { Sprite } from "./sprite.ts";
import { Canvas } from "./canvas.ts";

// Test for setting the position
Deno.test("Test SetPosition", () => {
  const sprite = new Sprite();
  const x = 10;
  const y = 20;
  sprite.setPosition(x, y);
  assertEquals(sprite.x, x);
  assertEquals(sprite.y, y);
});

// Test for setting the size
Deno.test("Test SetSize", () => {
  const sprite = new Sprite();
  const width = 100;
  const height = 200;
  sprite.setSize(width, height);
  assertEquals(sprite.width, width);
  assertEquals(sprite.height, height);
});

// Test for setting an image
Deno.test("Test SetImage", async () => {
  const sprite = new Sprite();
  const imgBytes = createTestImage();

  await sprite.setImage(imgBytes);
  assertNotEquals(sprite.image, null);
});

// Test for moving the sprite
Deno.test("Test Move", () => {
  const sprite = new Sprite();
  sprite.setPosition(10, 20);
  sprite.move(5, -5);
  assertEquals(sprite.x, 15);
  assertEquals(sprite.y, 15);
});

// Test for resizing the sprite
Deno.test("Test Resize", () => {
  const sprite = new Sprite();
  sprite.setSize(100, 200);
  sprite.resize(50, 100);
  assertEquals(sprite.width, 50);
  assertEquals(sprite.height, 100);
});

// Helper function to create a simple red 1x1 pixel image for testing
function createTestImage(): Uint8Array {
  const canvas = Canvas.create(1, 1);
  const ctx = canvas.getContext("2d")!;
  ctx.fillStyle = "red";
  ctx.fillRect(0, 0, 1, 1);
  return canvas.toBuffer();
}
