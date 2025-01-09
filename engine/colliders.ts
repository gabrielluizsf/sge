export interface Collider {
  collidesWithBox(box: BoxCollider): boolean;
}

export class BoxCollider implements Collider {
  constructor(
    public x: number,
    public y: number,
    public width: number,
    public height: number,
  ) {}

  collidesWithBox(box: BoxCollider): boolean {
    return (
      this.x < box.x + box.width &&
      this.x + this.width > box.x &&
      this.y < box.y + box.height &&
      this.y + this.height > box.y
    );
  }
}

export class CircleCollider implements Collider {
  constructor(public x: number, public y: number, public radius: number) {}

  collidesWithBox(box: BoxCollider): boolean {
    // Encontra o ponto mais próximo do círculo dentro do box
    const closestX = Math.max(box.x, Math.min(this.x, box.x + box.width));
    const closestY = Math.max(box.y, Math.min(this.y, box.y + box.height));

    // Calcula a distância entre o ponto mais próximo e o centro do círculo
    const distanceX = this.x - closestX;
    const distanceY = this.y - closestY;

    // Se a distância for menor ou igual ao raio, há colisão
    return distanceX * distanceX + distanceY * distanceY <=
      this.radius * this.radius;
  }
}
