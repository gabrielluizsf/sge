import { Time } from "./time.ts";
import type { GameObject } from "./game_object.ts";
import type { KeyboardEvent, Animator } from "./changed_types.ts";

export interface Updatable {
  update(deltaTime: number, keyboard: KeyboardObserver): void;
  gameObject: GameObject;
}

class Engine {
  private updatables: Updatable[] = [];
  private time: Time = new Time();

  public addUpdatable(updatable: Updatable): void {
    if (!this.updatables.includes(updatable)) {
      this.updatables.push(updatable);
    }
  }

  public objects(): number {
    return this.updatables.length;
  }

  public getGameObject(objectName: string): GameObject | undefined {
    return this.updatables.find((item) => item.gameObject.name === objectName)
      ?.gameObject;
  }

  public removeUpdatable(updatable: Updatable): void {
    this.updatables = this.updatables.filter((item) => item !== updatable);
  }

  public start(): void {
    this.time = new Time();
    this.loop();
  }

  private loop(): void {
    const deltaTime = this.time.deltaTime();
    const keyboard = KeyboardObserver.getInstance();

    for (const updatable of this.updatables) {
      updatable.update(deltaTime, keyboard);
    }
    (self as unknown as Animator).requestAnimationFrame?.(this.loop.bind(this));
  }
}

type ScriptConstructor = (constructor: new () => Updatable) => void;

export function RegisterWithEngine(engine: Engine): ScriptConstructor {
  return function (constructor: new () => Updatable) {
    const instance = new constructor();
    engine.addUpdatable(instance);
    engine.start();
  };
}

export const engine: Engine = new Engine();

export class KeyboardObserver {
  private static instance: KeyboardObserver;
  public keyState: Map<string, boolean> = new Map();

  private constructor() {
    if (typeof self.window !== "undefined") {
      self.window.addEventListener("keydown", (e: Event) => {
        if (this.isKeyboardEvent(e)) {
          this.keyState.set(e.key, true);
        }
      });

      self.window.addEventListener("keyup", (e: Event) => {
        if (this.isKeyboardEvent(e)) {
          this.keyState.set(e.key, false);
        }
      });
    }
  }

  private isKeyboardEvent(e: Event): e is KeyboardEvent {
    return "key" in e;
  }

  public static getInstance(): KeyboardObserver {
    if (!KeyboardObserver.instance) {
      KeyboardObserver.instance = new KeyboardObserver();
    }
    return KeyboardObserver.instance;
  }

  public isKeyPressed(key: string): boolean {
    return this.keyState.get(key) || false;
  }
}
