// internal/cpu: no instruction set extensions to report.
import * as $ from "../../../rt/mod.ts";

export const CacheLineSize = 64;

export class CacheLinePad extends $.Native {}
$.nativeType("cpu.CacheLinePad", CacheLinePad);

const none: any = new Proxy({}, { get: () => false });
export const X86 = none, ARM = none, ARM64 = none, Loong64 = none, MIPS64X = none, PPC64 = none, S390X = none, RISCV64 = none;

export function Initialize(_env: string): void {}

export function Name(): string {
  return "";
}
