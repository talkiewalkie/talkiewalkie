declare global {
  interface Object {
    let<T, K>(
      this: NonNullable<T>,
      functionToApply: (arg: NonNullable<T>) => K,
    ): K;
  }
}

const defineObjectProperty = <T extends keyof Object>(
  key: T,
  value: Object[T],
) => {
  Object.defineProperty(Object.prototype, key, { value });
};

defineObjectProperty("let", function (this, functionToApply) {
  // @ts-ignore
  return functionToApply(this.valueOf());
});

export {};
