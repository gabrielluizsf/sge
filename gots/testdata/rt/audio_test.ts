// The audio bridge asks SDL2 for float32 stereo and refuses a device that
// would play anything else, which would read the bytes it is given wrong.

import { accepted, spec } from "../../runtime/rt/audio.ts";

function assert(ok: boolean, what: string) {
  if (!ok) throw new Error(what);
}

Deno.test("the spec asks for float32 stereo at the rate", () => {
  const v = new DataView(spec(48000, 512).buffer);
  assert(v.getInt32(0, true) === 48000 && v.getUint16(4, true) === 0x8120 && v.getUint8(6) === 2, "spec");
});

Deno.test("a device of float32 stereo is accepted at its rate", () => {
  assert(accepted(spec(44100, 512)) === 44100, "accepted");
});

Deno.test("a device of another format, channels or rate is refused", () => {
  const s16 = spec(44100, 512);
  new DataView(s16.buffer).setUint16(4, 0x8010, true);
  const mono = spec(44100, 512);
  new DataView(mono.buffer).setUint8(6, 1);
  for (const [name, s] of [["s16", s16], ["mono", mono], ["rate 0", spec(0, 512)], ["rate huge", spec(1 << 30, 512)]] as const) {
    assert(accepted(s) === 0, name + " accepted");
  }
});
