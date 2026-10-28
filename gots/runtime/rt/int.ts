// Go's int, int64 and uint are kept in JavaScript numbers, exact up to 2^53.
// Arithmetic on them needs no wrapping, but JavaScript's bitwise operators
// cut to 32 bits, so the ones here keep the high bits.

const two32 = 4294967296;

function hi(x: number): number {
  return Math.floor(x / two32);
}

function lo(x: number): number {
  return x >>> 0;
}

function join(h: number, l: number): number {
  return (h | 0) * two32 + (l >>> 0);
}

function small(x: number): boolean {
  return x >= -2147483648 && x <= 2147483647;
}

export function and(a: number, b: number): number {
  if (small(a) && small(b)) return a & b;
  return join(hi(a) & hi(b), lo(a) & lo(b));
}

export function or(a: number, b: number): number {
  if (small(a) && small(b)) return a | b;
  return join(hi(a) | hi(b), lo(a) | lo(b));
}

export function xor(a: number, b: number): number {
  if (small(a) && small(b)) return a ^ b;
  return join(hi(a) ^ hi(b), lo(a) ^ lo(b));
}

export function andNot(a: number, b: number): number {
  if (small(a) && small(b)) return a & ~b;
  return join(hi(a) & ~hi(b), lo(a) & ~lo(b));
}

export function shl(a: number, n: number): number {
  if (n >= 64) return 0;
  return a * 2 ** n;
}

export function shr(a: number, n: number): number {
  if (n >= 64) return a < 0 ? -1 : 0;
  return Math.floor(a / 2 ** n);
}

// shl32 and shr32 are shifts of a 32-bit value by a count that may reach
// past its width, where Go gives zero and JavaScript wraps the count.
export function shl32(a: number, n: number): number {
  return n >= 32 ? 0 : a << n;
}

export function shr32(a: number, n: number): number {
  return n >= 32 ? (a < 0 ? -1 : 0) : a >> n;
}

export function shrU32(a: number, n: number): number {
  return n >= 32 ? 0 : a >>> n;
}

// div is integer division, panicking on zero as Go does.
export function div(a: number, b: number): number {
  if (b === 0) throw new RangeError("integer divide by zero");
  return Math.trunc(a / b);
}

export function rem(a: number, b: number): number {
  if (b === 0) throw new RangeError("integer divide by zero");
  return a % b;
}

// trunc is float-to-integer conversion: towards zero, NaN as zero.
export function trunc(x: number): number {
  return Number.isNaN(x) ? 0 : Math.trunc(x);
}

// u64 wraps a bigint to uint64, i64 to int64.
export function u64(x: bigint): bigint {
  return BigInt.asUintN(64, x);
}

export function divU64(a: bigint, b: bigint): bigint {
  if (b === 0n) throw new RangeError("integer divide by zero");
  return a / b;
}

export function shlU64(a: bigint, n: number): bigint {
  return n >= 64 ? 0n : BigInt.asUintN(64, a << BigInt(n));
}

export function shrU64(a: bigint, n: number): bigint {
  return n >= 64 ? 0n : a >> BigInt(n);
}

// toU64 converts an integer number to uint64, wrapping negatives.
export function toU64(x: number): bigint {
  return BigInt.asUintN(64, BigInt(Math.trunc(x)));
}

// fromU64 converts a uint64 to a number-backed integer, as int(x) does.
export function fromU64(x: bigint): number {
  return Number(BigInt.asIntN(64, x));
}
