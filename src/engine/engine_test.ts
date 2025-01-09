import { assertEquals } from "@std/assert/equals";
import {
  engine,
  KeyboardObserver,
  RegisterWithEngine,
  type Updatable,
} from "./engine.ts";
import { type Transform, Vector3 } from "./transform.ts";
import { GameObject } from "./game_object.ts";
import { BoxCollider, CircleCollider } from "./colliders.ts";
import { assertFalse } from "@std/assert/false";

@RegisterWithEngine(engine)
export class Player implements Updatable {
  gameObject: GameObject = new GameObject("player");
  private transform: Transform = this.gameObject.transform;
  isFoundEnemy: boolean = false;
  isDead: boolean = false;

  constructor() {
    this.transform.position = new Vector3(0, 0);
  }

  update(deltaTime: number, keyboard: KeyboardObserver): void {
    console.log(
      `Player position updated to (${this.transform.position.x}, ${this.transform.position.y}) delta time: ${deltaTime}`,
    );
    const colliderWidth = this.gameObject.sprite.width;
    const colliderHeight = this.gameObject.sprite.height;
    const collider = new BoxCollider(
      this.transform.position.x,
      this.transform.position.y,
      colliderWidth,
      colliderHeight,
    );
    this.gameObject.addCollider(collider);
    const hit = engine.getGameObject("enemy")?.checkCollision(this.gameObject);
    if (hit) {
      this.isFoundEnemy = true;
      this.isDead = true;
    }
    const isHelped = hit && keyboard.isKeyPressed("h");
    if (isHelped) {
      this.isDead = false;
    }
    if (keyboard.isKeyPressed("w")) {
      this.move(0, 1);
    }
    if (keyboard.isKeyPressed("s")) {
      this.move(0, -1);
    }
    if (keyboard.isKeyPressed("a")) {
      this.move(-1, 0);
    }
    if (keyboard.isKeyPressed("d")) {
      this.move(1, 0);
    }
  }

  private move(x: number, y: number): void {
    this.transform.position.x += x;
    this.transform.position.y += y;
  }

  getPosition() {
    return this.transform.position;
  }
}

@RegisterWithEngine(engine)
export class Enemy implements Updatable {
  gameObject: GameObject = new GameObject("enemy");
  update(deltaTime: number, _: KeyboardObserver): void {
    this.gameObject.transform.position = new Vector3(0, 1);
    this.gameObject.addCollider(new CircleCollider(0, 0, 10));
    console.log(
      `Enemy position updated to (${this.gameObject.transform.position.x}, ${this.gameObject.transform.position.y}) delta time: ${deltaTime}`,
    );
  }
}

Deno.test("Player should move up when 'w' is pressed", () => {
  const enemy = new Enemy();
  let player = new Player();
  const keyboard = KeyboardObserver.getInstance();
  keyboard.keyState.set("w", true);
  enemy.update(1, keyboard);
  player.update(1, keyboard);
  assertEquals(player.getPosition(), new Vector3(0, 1));
  assertFalse(!(player.isDead));
  assertFalse(!(player.isFoundEnemy));
  player = new Player();
  assertEquals(engine.objects(), 2);
  assertFalse(player.isDead);
  keyboard.keyState.set("w", true);
  keyboard.keyState.set("h", true);
  player.update(1, keyboard);
  assertFalse(player.isDead);
  assertFalse(!(player.isFoundEnemy));
});
