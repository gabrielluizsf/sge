// animation_test.ts
import {
  assertEquals,
  assertExists,
  assertThrows,
} from "@std/assert";

import {
  Animation,
  type AnimationCondition,
  Animator,
  ErrInvalidImages,
} from "./animation.ts";
import { GameObject } from "./game_object.ts";
import { Canvas, type CanvasContext2D } from "./canvas.ts";

class TestFiles {
     static  async generateImages(): Promise<CanvasContext2D[]> {
        const img1 = await Canvas.loadImage({ imgBytes: createTestImage() });
        const img2 = await Canvas.loadImage({ imgBytes: createTestImage() });
        return [img1.getContext("2d"), img2.getContext("2d")];
    }

    static deleteFiles(...files: CanvasContext2D[]): void {
        for (const file of files) {
            file.deleteFile();
        }
    }
}

function createTestImage(): Uint8Array {
    const canvas = Canvas.create(1, 1);
    const ctx = canvas.getContext("2d");
    ctx.fillStyle = "red";
    ctx.fillRect(0, 0, 1, 1);
    return canvas.toBuffer();
  }


Deno.test("Animation - should create animation with valid parameters", async() => {
  const condition: AnimationCondition = { name: "test", condition: true };
  const files = await TestFiles.generateImages();
  assertEquals(files.length, 2);
  const [firstImage, secondImage] = files;
  let err = firstImage.setFilePath("image1.png").createFile();
  assertEquals(err, null);
  err = secondImage.setFilePath("image2.png").createFile();
  assertEquals(err, null);

  
  const images = ["image1.png", "image2.png"];
  const animation = new Animation(condition, images, 1, "testAnimation");

  assertEquals(animation.name, "testAnimation");
  assertEquals(animation.images, images);
  assertEquals(animation.controller, condition);
  assertEquals(animation.isFinished, false);
  TestFiles.deleteFiles(firstImage, secondImage);
});

Deno.test("Animation - should throw error when images array is empty", () => {
  const condition: AnimationCondition = { name: "test", condition: true };

  assertThrows(
    () => new Animation(condition, [], 1),
    ErrInvalidImages,
    "O array de imagens não pode estar vazio.",
  );
});

Deno.test("Animation - should use default name when name is not provided", () => {
  const condition: AnimationCondition = { name: "test", condition: true };
  const images = ["image1.png"];
  const animation = new Animation(condition, images, 1);

  assertEquals(animation.name, "FooAnimation");
});

Deno.test("Animation - should stop animation correctly", () => {
  const condition: AnimationCondition = { name: "test", condition: true };
  const images = ["image1.png"];
  const animation = new Animation(condition, images, 1);

  animation.start();
  animation.stop();

  assertEquals(animation.isFinished, true);
});


Deno.test("Animator - should create animator with game object", () => {
  const gameObject = new GameObject("testObject");
  const animator = new Animator(gameObject);

  assertExists(animator);
});

Deno.test("Animator - should add animation successfully", async () => {
  const [firstImage] = await TestFiles.generateImages();
  const gameObject = new GameObject("testObject");
  const animator = new Animator(gameObject);
  const condition: AnimationCondition = { name: "test", condition: true };
  const animation = new Animation(
    condition,
    ["image1.png"],
    1,
    "testAnimation",
  );

  const result = animator.addAnimation(animation);
  assertEquals(result, true);
  TestFiles.deleteFiles(firstImage);
});

Deno.test("Animator - should not add duplicate animation", async() => {
  const [firstImage] = await TestFiles.generateImages();
  const gameObject = new GameObject("testObject");
  const animator = new Animator(gameObject);
  const condition: AnimationCondition = { name: "test", condition: true };
  const animation = new Animation(
    condition,
    ["image1.png"],
    1,
    "testAnimation",
  );

  animator.addAnimation(animation);
  const result = animator.addAnimation(animation);
  assertEquals(result, false);
  TestFiles.deleteFiles(firstImage);
});

Deno.test("Animator - should remove animation successfully", async() => {
  const [firstImage] = await TestFiles.generateImages();
  const gameObject = new GameObject("testObject");
  const animator = new Animator(gameObject);
  const condition: AnimationCondition = { name: "test", condition: true };
  const animation = new Animation(
    condition,
    ["image1.png"],
    1,
    "testAnimation",
  );

  animator.addAnimation(animation);
  const result = animator.removeAnimation("testAnimation");
  assertEquals(result, true);
  TestFiles.deleteFiles(firstImage);
});

Deno.test("Animator - should return false when removing non-existent animation", () => {
  const gameObject = new GameObject("testObject");
  const animator = new Animator(gameObject);

  const result = animator.removeAnimation("nonExistentAnimation");
  assertEquals(result, false);
});


Deno.test("Animator - should not play animation if condition is false", async () => {
  const gameObject = new GameObject("testObject");
  const animator = new Animator(gameObject);
  const condition: AnimationCondition = { name: "test", condition: false };
  const animation = new Animation(
    condition,
    ["image1.png"],
    1,
    "testAnimation",
  );

  animator.addAnimation(animation);
  await animator.play("testAnimation");

  assertEquals(animation.isFinished, false);
});
