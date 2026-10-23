// Go's predeclared error interface.
export interface error {
  Error(): string;
}

// minOf and maxOf are min and max for strings and bigints, which Math's
// cannot compare.
export function minOf<T>(...xs: T[]): T {
  let m = xs[0];
  for (const x of xs) if (x < m) m = x;
  return m;
}

export function maxOf<T>(...xs: T[]): T {
  let m = xs[0];
  for (const x of xs) if (x > m) m = x;
  return m;
}
