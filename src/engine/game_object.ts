import { Transform, type Vector3 } from "./transform.ts";
import { BoxCollider, CircleCollider, type Collider } from "./colliders.ts";
import { Sprite } from "./sprite.ts";

// GameObject class represents an object in the game
export class GameObject {
    name: string;
    transform: Transform;
    collider: Collider = new BoxCollider(0, 0, 0, 0);
    sprite: Sprite;
  
    constructor(name: string) {
      this.name = name;
      this.transform = new Transform();
      this.sprite = new Sprite();
    }
  
    // Method to change the sprite's image
    async changeSprite(imgBytes: Uint8Array): Promise<void> {
      const sprite = new Sprite();
      await sprite.setImage(imgBytes);
      this.sprite = sprite;
    }
  
    // Method to add a collider to the GameObject
    addCollider(collider: Collider): void {
      this.collider = collider;
    }
  
    // Method to update the GameObject's position (move, rotate, scale)
    updatePosition(delta: Vector3): void {
      this.transform.translate(delta);
    }
  
    // Method to check for collision with another GameObject
    checkCollision(other: GameObject): boolean {
      let collider = new BoxCollider(0, 0, 0, 0);
      if (other.collider instanceof BoxCollider) {
        collider = other.collider;
      } else if (other.collider instanceof CircleCollider) {
        collider = this.convertCircleToBox(other.collider);
      }
      return this.collider.collidesWithBox(collider);
    }
  
    // Convert CircleCollider to BoxCollider
    convertCircleToBox(c: CircleCollider): BoxCollider {
      return new BoxCollider(
        c.x - c.radius, 
        c.y - c.radius, 
        2 * c.radius, 
        2 * c.radius
      );
    }
  }