// Go's math package over JavaScript's Math, with Go's answers where the two
// disagree (Round, Signbit, the bit conversions).
export const Pi = Math.PI;
export const E = Math.E;
export const Sqrt2 = Math.SQRT2;
export const Ln2 = Math.LN2;
export const MaxFloat64 = Number.MAX_VALUE;
export const SmallestNonzeroFloat64 = Number.MIN_VALUE;
export const MaxInt = Number.MAX_SAFE_INTEGER;
export const MinInt = Number.MIN_SAFE_INTEGER;

export const Abs = Math.abs;
export const Sqrt = Math.sqrt;
export const Cbrt = Math.cbrt;
export const Sin = Math.sin;
export const Cos = Math.cos;
export const Tan = Math.tan;
export const Asin = Math.asin;
export const Acos = Math.acos;
export const Atan = Math.atan;
export const Atan2 = Math.atan2;
export const Sinh = Math.sinh;
export const Cosh = Math.cosh;
export const Tanh = Math.tanh;
export const Exp = Math.exp;
export const Expm1 = Math.expm1;
export const Log = Math.log;
export const Log2 = Math.log2;
export const Log10 = Math.log10;
export const Log1p = Math.log1p;
export const Pow = Math.pow;
export const Floor = Math.floor;
export const Ceil = Math.ceil;
export const Trunc = Math.trunc;
export const Max = Math.max;
export const Min = Math.min;
export const Hypot = Math.hypot;

export function Sincos(x: number): [number, number] {
  return [Math.sin(x), Math.cos(x)];
}

// Round rounds half away from zero, where Math.round rounds half up.
export function Round(x: number): number {
  return x < 0 ? -Math.round(-x) : Math.round(x);
}

export function RoundToEven(x: number): number {
  const r = Math.round(x);
  return Math.abs(x % 1) === 0.5 && r % 2 !== 0 ? r - 1 : r;
}

export function Mod(x: number, y: number): number {
  return x % y;
}

export function Modf(x: number): [number, number] {
  const i = Math.trunc(x);
  return [i, x - i];
}

export function Inf(sign: number): number {
  return sign >= 0 ? Infinity : -Infinity;
}

export function NaN(): number {
  return globalThis.NaN;
}

export function IsNaN(x: number): boolean {
  return x !== x;
}

export function IsInf(x: number, sign: number): boolean {
  return (sign >= 0 && x === Infinity) || (sign <= 0 && x === -Infinity);
}

export function Signbit(x: number): boolean {
  return x < 0 || Object.is(x, -0);
}

export function Copysign(x: number, s: number): number {
  return Signbit(x) === Signbit(s) ? x : -x;
}

export function Dim(x: number, y: number): number {
  return x > y ? x - y : 0;
}

const view = new DataView(new ArrayBuffer(8));

export function Float64bits(f: number): bigint {
  view.setFloat64(0, f);
  return view.getBigUint64(0);
}

export function Float64frombits(b: bigint): number {
  view.setBigUint64(0, b);
  return view.getFloat64(0);
}

export function Float32bits(f: number): number {
  view.setFloat32(0, f);
  return view.getUint32(0);
}

export function Float32frombits(b: number): number {
  view.setUint32(0, b);
  return view.getFloat32(0);
}
