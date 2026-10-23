// Float formatting with Go's rules, which differ from JavaScript's in when
// they switch to exponent form and in how exponents are written.

// digits is the shortest decimal that reads back as f: its digits, and the
// position of the decimal point relative to them.
function shortest(f: number, bits: number): { d: string; dp: number } {
  let s: string;
  if (bits === 32) {
    s = "";
    for (let p = 1; p <= 9; p++) {
      s = f.toExponential(p - 1);
      if (Math.fround(parseFloat(s)) === f) break;
    }
  } else {
    s = f.toExponential();
  }
  return split(s);
}

// split takes JavaScript's d.ddde±x apart.
function split(s: string): { d: string; dp: number } {
  const e = s.indexOf("e");
  const d = s.slice(0, e).replace(".", "").replace("-", "");
  return { d, dp: parseInt(s.slice(e + 1), 10) + 1 };
}

// halfEven fixes JavaScript's rounding of an exact tie, which goes up, to
// Go's, which goes to the even digit. exact is the value's digits written
// far past the cut, rounded is JavaScript's result with the same digits up
// to it.
function tieDown(exact: string, keep: number): boolean {
  const digits = exact.replace(/[-.]/g, "").replace(/e.*$/, "");
  if (digits[keep] !== "5" || !/^0*$/.test(digits.slice(keep + 1))) return false;
  const last = keep > 0 ? digits.charCodeAt(keep - 1) - 48 : 0;
  return last % 2 === 0;
}

function rounded(f: number, prec: number): { d: string; dp: number } {
  const p = Math.min(Math.max(prec - 1, 0), 100);
  if (p < 100 && tieDown(Math.abs(f).toExponential(100), p + 1)) {
    const exact = split(Math.abs(f).toExponential(100));
    return { d: exact.d.slice(0, p + 1), dp: exact.dp };
  }
  return split(Math.abs(f).toExponential(p));
}

// toFixedGo is toFixed rounding ties to even.
function toFixedGo(a: number, prec: number): string {
  if (prec < 100) {
    const exact = a.toFixed(100);
    const point = exact.indexOf(".");
    if (tieDown(exact, point + prec)) {
      const cut = exact.slice(0, point + 1 + prec);
      return cut.endsWith(".") ? cut.slice(0, -1) : cut;
    }
  }
  return a.toFixed(prec);
}

function fixedForm(d: string, dp: number, prec: number): string {
  let int: string, frac: string;
  if (dp <= 0) {
    int = "0";
    frac = "0".repeat(-dp) + d;
  } else if (dp >= d.length) {
    int = d + "0".repeat(dp - d.length);
    frac = "";
  } else {
    int = d.slice(0, dp);
    frac = d.slice(dp);
  }
  if (prec >= 0) frac = frac.padEnd(prec, "0").slice(0, prec);
  return frac.length > 0 ? int + "." + frac : int;
}

function trimZeros(d: string): string {
  return d.replace(/0+$/, "") || "0";
}

// eDigits prints d as d.ddd with prec digits after the point and an
// exponent of at least two digits, as strconv's %e does.
function eForm(d: string, dp: number, prec: number, upper: boolean): string {
  let m = d[0];
  if (prec > 0) m += "." + d.slice(1).padEnd(prec, "0").slice(0, prec);
  const x = d === "0" ? 0 : dp - 1;
  const ax = Math.abs(x);
  return m + (upper ? "E" : "e") + (x < 0 ? "-" : "+") + (ax < 10 ? "0" + ax : globalThis.String(ax));
}

// formatFloat is strconv.FormatFloat for the formats fmt uses: 'e', 'E', 'f',
// 'g' and 'G', with prec -1 for the shortest form that reads back exactly.
export function formatFloat(f: number, fmt: string, prec: number, bits: number): string {
  if (f !== f) return "NaN";
  if (f === Infinity) return "+Inf";
  if (f === -Infinity) return "-Inf";
  const neg = f < 0 || Object.is(f, -0);
  const a = Math.abs(f);
  const shortestForm = prec < 0;
  let out: string;
  if (fmt === "f" || fmt === "F") {
    if (shortestForm) {
      const { d, dp } = a === 0 ? { d: "0", dp: 1 } : shortest(a, bits);
      out = fixedForm(d, dp, Math.max(d.length - dp, 0));
    } else if (a < 1e21) {
      out = toFixedGo(a, Math.min(prec, 100));
    } else {
      const { d, dp } = shortest(a, bits);
      out = fixedForm(d, dp, prec);
    }
  } else if (fmt === "e" || fmt === "E") {
    let { d, dp } = a === 0 ? { d: "0", dp: 1 } : shortestForm ? shortest(a, bits) : rounded(a, prec + 1);
    if (shortestForm) prec = d.length - 1;
    out = eForm(d, dp, prec, fmt === "E");
  } else {
    let d: string, dp: number;
    if (a === 0) {
      d = "0", dp = 1;
    } else {
      ({ d, dp } = shortestForm ? shortest(a, bits) : rounded(a, prec === 0 ? 1 : prec));
    }
    d = trimZeros(d);
    const nd = d.length;
    if (shortestForm) prec = nd;
    else if (prec === 0) prec = 1;
    let eprec = prec;
    if (eprec > nd && nd >= dp) eprec = nd;
    if (shortestForm) eprec = 6;
    const x = dp - 1;
    if (x < -4 || x >= eprec) {
      out = eForm(d, dp, Math.min(prec, nd) - 1, fmt === "G");
    } else {
      out = fixedForm(d, dp, Math.max((prec > dp ? nd : prec) - dp, 0));
    }
  }
  return neg ? "-" + out : out;
}
