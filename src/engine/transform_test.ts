import { assertEquals, assertFalse } from "@std/assert";
import { Quaternion, Transform, Vector3 } from "./transform.ts"; 

Deno.test("TestTranslate", () => {
    const t1 = new Transform();
  
    t1.translate(new Vector3(1, 2, 3));
  
    const expected = new Vector3(1, 2, 3);
  
    assertEquals(t1.position, expected);
  });
  
  Deno.test("TestRotate", () => {
    const t1 = new Transform();
  
    const axis = new Vector3(0, 1, 0);
    const angle = Math.PI / 2;
    t1.rotate(angle, axis);
  
    const sinHalfAngle = Math.sin(angle / 2);
    const cosHalfAngle = Math.cos(angle / 2);
    const expectedRotation = new Quaternion(
      axis.x * sinHalfAngle,
      axis.y * sinHalfAngle,
      axis.z * sinHalfAngle,
      cosHalfAngle,
    );
  
    assertEquals(t1.rotation, expectedRotation);
  
    t1.rotation = Quaternion.identity;
    assertEquals(t1.rotation, Quaternion.identity);
  });
  
  Deno.test("TestScaleBy", () => {
    const t1 = new Transform();
  
    t1.scaleBy(new Vector3(2, 2, 2));
  
    const expected = new Vector3(0, 0, 0); // Since initial scale is (0, 0, 0), multiplying by any factor keeps it at (0, 0, 0)
    assertEquals(t1.scale, expected);
  });
  
  Deno.test("TestIsMoving", () => {
    const t1 = new Transform();
  
    assertFalse(t1.position.sqrMagnitude() > 0);
  
    t1.translate(new Vector3(1, 2, 3));
    const isMoving = t1.position.sqrMagnitude() > 0;
    assertFalse(!isMoving);
  });
  