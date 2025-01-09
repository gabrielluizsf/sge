import { assertInstanceOf, assertEquals, assertExists } from "@std/assert";
import { Canvas, ErrInvalidImageOptions } from "./canvas.ts"; 

Deno.test("Canvas.create should create a canvas with the specified width and height", () => {
  const canvas = Canvas.create(200, 300);
  const context = canvas.getContext("2d");
  assertExists(context)
  assertEquals(canvas.width, 200)
  assertEquals(canvas.height, 300)
});

Deno.test("Canvas.loadImage should throw an error for invalid options", async () => {
  try {
    await Canvas.loadImage({});
    throw new Error("Expected ErrInvalidImageOptions to be thrown");
  } catch (error) {
    assertInstanceOf(error, ErrInvalidImageOptions);
  }
});
