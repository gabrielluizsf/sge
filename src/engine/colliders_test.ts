import { assertEquals } from "@std/assert";
import { BoxCollider, CircleCollider } from "./colliders.ts"; 

Deno.test("BoxCollider: collidesWithBox - overlapping boxes", () => {
  const box1 = new BoxCollider(0, 0, 10, 10);
  const box2 = new BoxCollider(5, 5, 10, 10);

  assertEquals(box1.collidesWithBox(box2), true);
  assertEquals(box2.collidesWithBox(box1), true);
});

Deno.test("BoxCollider: collidesWithBox - non-overlapping boxes", () => {
  const box1 = new BoxCollider(0, 0, 10, 10);
  const box2 = new BoxCollider(20, 20, 10, 10);

  assertEquals(box1.collidesWithBox(box2), false);
  assertEquals(box2.collidesWithBox(box1), false);
});

Deno.test("BoxCollider: collidesWithBox - touching edges", () => {
  const box1 = new BoxCollider(0, 0, 10, 10);
  const box2 = new BoxCollider(10, 0, 10, 10);

  assertEquals(box1.collidesWithBox(box2), false);
  assertEquals(box2.collidesWithBox(box1), false);
});

Deno.test("CircleCollider: collidesWithBox - circle inside box", () => {
  const box = new BoxCollider(0, 0, 10, 10);
  const circle = new CircleCollider(5, 5, 3);

  assertEquals(circle.collidesWithBox(box), true);
});

Deno.test("CircleCollider: collidesWithBox - circle overlapping box edge", () => {
  const box = new BoxCollider(0, 0, 10, 10);
  const circle = new CircleCollider(9, 5, 3);

  assertEquals(circle.collidesWithBox(box), true);
});

Deno.test("CircleCollider: collidesWithBox - circle outside box", () => {
  const box = new BoxCollider(0, 0, 10, 10);
  const circle = new CircleCollider(15, 15, 3);

  assertEquals(circle.collidesWithBox(box), false);
});

Deno.test("CircleCollider: collidesWithBox - circle touching box edge", () => {
  const box = new BoxCollider(0, 0, 10, 10);
  const circle = new CircleCollider(10, 5, 1);

  assertEquals(circle.collidesWithBox(box), true);
});
