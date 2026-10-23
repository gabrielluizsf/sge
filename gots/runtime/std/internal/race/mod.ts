// internal/race: there is no race detector.
const nothing = (..._a: any[]): void => {};

export const Acquire = nothing, Release = nothing, ReleaseMerge = nothing, Disable = nothing, Enable = nothing;
export const Read = nothing, ReadPC = nothing, ReadObjectPC = nothing, ReadRange = nothing;
export const Write = nothing, WritePC = nothing, WriteObjectPC = nothing, WriteRange = nothing;

export function Errors(): number {
  return 0;
}
