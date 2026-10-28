// Go's iter package. Seq and Seq2 are function types; Pull runs the whole
// sequence up front, as there are no coroutines to pause it with.
import * as $ from "../../rt/mod.ts";

export type Seq<V> = ((yield_: (v: V) => boolean) => void) | null;
export type Seq2<K, V> = ((yield_: (k: K, v: V) => boolean) => void) | null;
export const Seq$type = $.named("iter.Seq", $.funcType("func(func(V) bool)"), null);
export const Seq2$type = $.named("iter.Seq2", $.funcType("func(func(K, V) bool)"), null);

export function Pull<V>(seq: Seq<V>): [() => [V | null, boolean], () => void] {
  const vals: V[] = [];
  seq!((v: V) => {
    vals.push(v);
    return true;
  });
  let i = 0;
  return [() => (i < vals.length ? [vals[i++], true] : [null, false]), () => { i = vals.length; }];
}

export function Pull2<K, V>(seq: Seq2<K, V>): [() => [K | null, V | null, boolean], () => void] {
  const vals: [K, V][] = [];
  seq!((k: K, v: V) => {
    vals.push([k, v]);
    return true;
  });
  let i = 0;
  return [
    () => (i < vals.length ? [vals[i][0], vals[i++][1], true] : [null, null, false]),
    () => { i = vals.length; },
  ];
}
