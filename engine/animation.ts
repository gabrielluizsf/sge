import { GameObject } from "./game_object.ts";
import { convertImageToUint8Array } from "./canvas.ts";

export type AnimationCondition = {
  name: string;
  condition: boolean;
};

export class ErrInvalidImages implements Error {
  name: string = "ErrInvalidImages";
  message: string = "O array de imagens não pode estar vazio.";
  stack?: string | undefined;
  cause?: unknown;
  constructor(cause?: unknown) {
    this.name = "ErrInvalidImageOptions";
    this.message = "Invalid image options";
    this.cause = cause;
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ErrInvalidImages);
    }
  }
}

export class Animation {
  private currentIndex: number = 0;
  private intervalId: number | undefined;
  private interval: number;

  public name: string = "FooAnimation";
  public controller: AnimationCondition;
  public images: string[];
  public isFinished: boolean = false;

  constructor(
    animationCondition: AnimationCondition,
    images: string[],
    interval: number,
    name?: string,
  ) {
    if (images.length === 0) {
      throw new ErrInvalidImages();
    }

    this.images = images;
    this.interval = interval;
    this.name = name ?? this.name;
    this.controller = { ...animationCondition };
  }

  start(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }

    this.intervalId = setInterval(() => {
      this.changeImage();
    }, this.animationInterval());
  }

  private animationInterval(): number {
    return this.interval * 1000;
  }

  stop(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
      this.isFinished = true;
      this.intervalId = undefined;
    }
  }

  private changeImage(): void {
    if (this.images.length === 0) return;
    console.log(`Trocando para a imagem: ${this.images[this.currentIndex]}`);
    this.currentIndex = (this.currentIndex + 1) % this.images.length;
  }
}

export class Animator {
  private animations: Animation[] = [];
  private currentAnimation: Animation | undefined;
  private gameObject: GameObject = new GameObject("foo");
  private initialSprite: Uint8Array = new Uint8Array();

  constructor(gameObject: GameObject) {
    this.gameObject = gameObject;
    this.initialSprite = convertImageToUint8Array(gameObject.sprite.image);
  }

  #restoreSprite(): void {
    this.gameObject.changeSprite(this.initialSprite);
  }

  addAnimation(animation: Animation): boolean {
    if (this.animations.find((item) => item.name === animation.name)) {
      console.warn(`Animação com o nome "${animation.name}" já existe.`);
      return false;
    }
    this.animations.push(animation);
    return true;
  }

  removeAnimation(animationName: string): boolean {
    const index = this.animations.findIndex((item) =>
      item.name === animationName
    );
    if (index !== -1) {
      this.animations[index].stop();
      this.animations.splice(index, 1);
      return true;
    }
    console.warn(`Animação "${animationName}" não encontrada.`);
    return false;
  }

  stopAll(): void {
    this.animations.filter((anim) => anim.name !== this.currentAnimation?.name)
      .forEach((item) => item.stop());
  }

  async play(animationName: string): Promise<void> {
    const animation = this.animations.find((item) =>
      item.name === animationName
    );
    if (animation && animation.controller.condition) {
      if (this.currentAnimation?.name === animationName) {
        console.info(`Animação "${animationName}" já está em execução.`);
        return;
      }

      this.currentAnimation = animation;
      this.stopAll();
      this.currentAnimation.start();

      try {
        const imageData = await this.loadImageAsUint8Array(
          this.currentAnimation.images[0],
        );
        this.gameObject.changeSprite(imageData);
      } catch (error) {
        console.error(`Erro ao carregar sprite: ${error}`);
      }

      this.currentAnimation = undefined;
    }
  }
  private async loadImageAsUint8Array(imagePath: string): Promise<Uint8Array> {
    if (imagePath.startsWith("http") || imagePath.startsWith("https")) {
      const response = await fetch(imagePath);
      if (!response.ok) {
        throw new Error(`Falha ao carregar a imagem: ${imagePath}`);
      }
      const blob = await response.blob();
      return new Uint8Array(await blob.arrayBuffer());
    } else {
      const fs = await import("node:fs/promises");
      const data = await fs.readFile(imagePath);
      return new Uint8Array(data.buffer);
    }
  }
}
