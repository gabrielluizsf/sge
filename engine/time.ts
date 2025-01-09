export class Time {
  private previousTime: number;

  constructor() {
    this.previousTime = Date.now();
  }

  public deltaTime(): number {
    const currentTime = Date.now();
    const deltaTime = (currentTime - this.previousTime) / 1000;
    this.previousTime = currentTime;
    return deltaTime;
  }
}
