// Complex numbers, as immutable { re, im } pairs.
export type Complex = { re: number; im: number };

export function add(a: Complex, b: Complex): Complex {
  return { re: a.re + b.re, im: a.im + b.im };
}

export function sub(a: Complex, b: Complex): Complex {
  return { re: a.re - b.re, im: a.im - b.im };
}

export function mul(a: Complex, b: Complex): Complex {
  return { re: a.re * b.re - a.im * b.im, im: a.re * b.im + a.im * b.re };
}

// div is Go's complex division, scaling to avoid overflow and handling the
// infinities and NaNs as the runtime does.
export function div(n: Complex, m: Complex): Complex {
  let e: number, f: number;
  if (Math.abs(m.re) >= Math.abs(m.im)) {
    const ratio = m.im / m.re;
    const denom = m.re + ratio * m.im;
    e = (n.re + n.im * ratio) / denom;
    f = (n.im - n.re * ratio) / denom;
  } else {
    const ratio = m.re / m.im;
    const denom = m.im + ratio * m.re;
    e = (n.re * ratio + n.im) / denom;
    f = (n.im * ratio - n.re) / denom;
  }
  if (Number.isNaN(e) && Number.isNaN(f)) {
    e = n.re / 0;
    f = n.im / 0;
  }
  return { re: e, im: f };
}

export function neg(a: Complex): Complex {
  return { re: -a.re, im: -a.im };
}

// c64 rounds both parts to float32, for complex64.
export function c64(a: Complex): Complex {
  return { re: Math.fround(a.re), im: Math.fround(a.im) };
}
