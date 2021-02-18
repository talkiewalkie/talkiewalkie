declare global {
  interface Array<T> {
    /**
     * Get first element of an array with nullable type
     */
    first(this: T[]): T | undefined;

    /**
     * Get last element of an array with nullable type
     */
    last(this: T[]): T | undefined;

    /**
     * Immutable version of {@link Array.sort}
     */
    immutableSort(this: T[], compareFn?: (a: T, b: T) => number): T[];

    /**
     * Shortcut for {@link Array.immutableSort}
     * - sortAsc((fn) => number) === immutableSort((a, b) => fn(a) - fn(b))
     * - sortAsc((fn) => string) === immutableSort((a, b) => fn(a).localCompare(fn(b)))
     */
    sortAsc(
      this: T[],
      getValue: ((el: T) => number) | ((el: T) => string),
    ): T[];

    /**
     * Shortcut for {@link Array.immutableSort}
     * - sortDesc((fn) => number) === immutableSort((a, b) => fn(b) - fn(a))
     * - sortDesc((fn) => string) === immutableSort((a, b) => fn(b).localCompare(fn(a)))
     */
    sortDesc(
      this: T[],
      getValue: ((el: T) => number) | ((el: T) => string),
    ): T[];

    /**
     * Shortcut for {@link Array.filter}.length
     */
    count(
      this: T[],
      predicate: (value: T, index: number, array: T[]) => unknown,
    ): number;

    /**
     * Typed shortcut of {@link Array.filter} for removing undefined and null
     * filterNotNull() === filter((el) => el !== undefined && el !== null)
     */
    filterNotNull(this: T[]): NonNullable<T>[];

    /**
     * Typed shortcut of {@link Array.map} and {@link Array.filterNotNull}
     * mapNotNull(predicate) === map(predicate).filterNotNull()
     */
    mapNotNull<U>(
      this: T[],
      callbackFn: (value: T, index: number, array: T[]) => U,
    ): NonNullable<U>[];

    /**
     * Shortcut for comparing {@link Array.length} to 0
     */
    isEmpty(this: T[]): boolean;

    /**
     * Returns a list containing only distinct (by reference) elements
     * Keeps the first in the array among duplicates
     */
    distinct(this: T[]): T[];

    /**
     * Returns a list containing only elements having distinct (by reference)
     * results for the selector function
     * Keeps the first in the array among duplicates
     */
    distinctBy(this: T[], selector: (value: T) => unknown): T[];

    /**
     * Groupby, the good stuff
     */
    groupBy(this: T[], selector: (value: T) => string): Record<string, T[]>;
  }
}

const defineArrayProperty = <K extends keyof Array<unknown>>(
  key: K,
  value: Array<unknown>[K],
) => {
  Object.defineProperty(Array.prototype, key, { value });
};

defineArrayProperty("first", function () {
  return this[0];
});

defineArrayProperty("last", function () {
  if (this.length > 0) return this[this.length - 1];
  return undefined;
});

defineArrayProperty("immutableSort", function (compareFn) {
  return this.slice().sort(compareFn);
});

defineArrayProperty("isEmpty", function () {
  return this.length === 0;
});

const isStringSort = <T>(
  getValue: (value: T) => string | number | Date,
  value: T,
): getValue is (value: T) => string => typeof getValue(value) === "string";

defineArrayProperty("sortAsc", function (getValue) {
  if (this.isEmpty()) return this;
  if (isStringSort(getValue, this[0])) {
    return this.immutableSort((a, b) => getValue(a).localeCompare(getValue(b)));
  }
  return this.immutableSort((a, b) => getValue(a) - getValue(b));
});

defineArrayProperty("sortDesc", function (getValue) {
  if (this.isEmpty()) return this;
  if (isStringSort(getValue, this[0])) {
    return this.immutableSort((a, b) => getValue(b).localeCompare(getValue(a)));
  }
  return this.immutableSort((a, b) => getValue(b) - getValue(a));
});

defineArrayProperty("count", function (predicate) {
  return this.filter(predicate).length;
});

defineArrayProperty("filterNotNull", function () {
  return this.filter((el) => el !== undefined && el !== null);
});

defineArrayProperty("mapNotNull", function (predicate) {
  return this.map(predicate).filterNotNull();
});

defineArrayProperty("distinctBy", function (selector) {
  const set = new Set();
  return this.filter((e) => {
    if (set.has(selector(e))) return false;
    set.add(selector(e));
    return true;
  });
});

defineArrayProperty("distinct", function () {
  return this.distinctBy((e) => e);
});

defineArrayProperty("groupBy", function (selector) {
  return this.reduce<Record<string, Array<unknown>>>((acc, it) => {
    acc[selector(it)] = selector(it) in acc ? [...acc[selector(it)], it] : [it];
    return acc;
  }, {});
});

export {};
