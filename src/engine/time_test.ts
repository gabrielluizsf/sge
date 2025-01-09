import { assertEquals } from "@std/assert";
import { Time } from "./time.ts"; 

Deno.test("Time.deltaTime calculates time difference correctly", () => {
  const mockDateNow = () => {
    let time = 1000;
    return {
      now: () => {
        const currentTime = time;
        time += 1000; 
        return currentTime;
      },
    };
  };

  const originalDateNow = Date.now;
  const mock = mockDateNow();
  Date.now = mock.now;

  try {
    const time = new Time();
    assertEquals(time.deltaTime(), 1); 
    assertEquals(time.deltaTime(), 1); 
  } finally {
    Date.now = originalDateNow; 
  }
});