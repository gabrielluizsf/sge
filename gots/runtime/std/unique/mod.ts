// Go's unique package: equal values made into handles share one, so
// handles compare by identity the way Go's compare by pointer.
import * as $ from "../../rt/mod.ts";

export class Handle<T> extends $.Native {
  value: T | null = null;
  Value(): T {
    return this.value as T;
  }
  override $eq(o: Handle<T>): boolean {
    return this.value === o.value || $.keyOf(this.value) === $.keyOf(o.value);
  }
  override $key(): string {
    return "h" + $.keyOf(this.value);
  }
}
$.nativeType("unique.Handle", Handle);

const handles = new Map<string, Handle<any>>();

export function Make<T>(value: T): Handle<T> {
  const k = $.keyOf(value);
  let h = handles.get(k);
  if (h === undefined) {
    h = new Handle<T>();
    h.value = value;
    handles.set(k, h);
  }
  return h;
}
