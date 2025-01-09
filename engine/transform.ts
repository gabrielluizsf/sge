export class Vector3 {
  x: number = 0;
  y: number = 0;
  z: number = 0;
  constructor(
    x?: number,
    y?: number,
    z?: number,
  ) {
    this.x = x ?? 0;
    this.y = y ?? 0;
    this.z = z ?? 0;
  }

  public sqrMagnitude(): number {
    return this.x * this.x + this.y * this.y + this.z * this.z;
  }
}

interface Q {
  x: number;
  y: number;
  z: number;
  w: number;
}

export class Quaternion implements Q {
  x: number = 0;
  y: number = 0;
  z: number = 0;
  w: number = 1;
  public static identity: Q = {
    x: 0,
    y: 0,
    z: 0,
    w: 1,
  };

  constructor(
    x: number,
    y: number,
    z: number,
    w: number,
  ) {
    this.x = x;
    this.y = y;
    this.z = z;
    this.w = w;
  }
}

export class Transform {
  position: Vector3 = new Vector3(0, 0, 0);
  rotation: Quaternion = Quaternion.identity;
  scale: Vector3 = new Vector3(0, 0, 0);
  public translate(delta: Vector3) {
    this.position.x += delta.x;
    this.position.y += delta.y;
    this.position.z += delta.z;
  }
  public rotate(angle: number, axis: Vector3) {
    const sin = Math.sin(angle / 2);
    const x = axis.x * sin;
    const y = axis.y * sin;
    const z = axis.z * sin;
    const w = Math.cos(angle / 2);
    this.rotation = new Quaternion(x, y, z, w);
  }
  public scaleBy(factor: Vector3) {
    this.scale.x *= factor.x;
    this.scale.y *= factor.y;
    this.scale.z *= factor.z;
  }
}
