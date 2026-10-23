// The parts of internal/abi the standard library uses outside the runtime.
export function NoEscape<T>(p: T): T {
  return p;
}

export function Escape<T>(x: T): T {
  return x;
}

export function FuncPCABI0(_f: any): number {
  return 0;
}

export function FuncPCABIInternal(_f: any): number {
  return 0;
}
