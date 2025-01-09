import { assertEquals, assertInstanceOf, assert } from "@std/assert";
import { GameObject } from "./game_object.ts";
import { BoxCollider, CircleCollider } from "./colliders.ts";
import { Sprite } from "./sprite.ts";

Deno.test("GameObject creation", () => {
  const name = "player";
  const player = new GameObject(name);

  assert(player);
  assertEquals(player.name, name);
  assert(player.transform);
  assert(player.sprite instanceof Sprite);
});

Deno.test("Add BoxCollider", () => {
  const mouse = new GameObject("mouse");
  mouse.addCollider(new BoxCollider(0, 0, 10, 10));

  assert(mouse.collider);
  assertInstanceOf(mouse.collider, BoxCollider);
});

Deno.test("Add CircleCollider", () => {
  const cat = new GameObject("cat");
  cat.addCollider(new CircleCollider(5, 5, 5));

  assert(cat.collider);
  assertInstanceOf(cat.collider, CircleCollider);
});

Deno.test("Check BoxCollider collision", () => {
  const cat = new GameObject("cat");
  cat.addCollider(new BoxCollider(0, 0, 10, 10));

  const dog = new GameObject("dog");
  dog.addCollider(new BoxCollider(5, 5, 10, 10));

  const collision = cat.checkCollision(dog);
  assertEquals(collision, true);
});

Deno.test("Check Circle and BoxCollider collision", () => {
  const cat = new GameObject("cat");
  cat.addCollider(new CircleCollider(5, 5, 5));

  const dog = new GameObject("dog");
  dog.addCollider(new BoxCollider(0, 0, 10, 10));

  const collision = cat.checkCollision(dog);
  assertEquals(collision, true);
});

Deno.test("No collision", () => {
  const cat = new GameObject("cat");
  cat.addCollider(new BoxCollider(0, 0, 10, 10));

  const dog = new GameObject("dog");
  dog.addCollider(new CircleCollider(20, 20, 5));

  const collision = cat.checkCollision(dog);
  assertEquals(collision, false);
});
