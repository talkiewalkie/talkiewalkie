declare global {
  interface String {
    addNonBrakingSpaces(this: string): string;
    upperFirst(this: string): string;
    truncate(this: string, maxLength: number): string;
    removeAccents(this: string): string;
    toSearchValue(this: string): string;
    fuzzyMatch(this: string, search: string): boolean;
  }
}

const defineStringProperty = <K extends keyof string>(
  key: K,
  value: string[K],
) => {
  Object.defineProperty(String.prototype, key, { value });
};

defineStringProperty("addNonBrakingSpaces", function () {
  return this.replaceAll(/ ([?|!])/gu, "\xa0$1");
});

defineStringProperty("upperFirst", function () {
  return this.length === 0 ? "" : this[0].toUpperCase() + this.slice(1);
});

defineStringProperty("truncate", function (maxLength) {
  const singleLine = this.replace(/\s+/gu, " ");
  return singleLine.length < maxLength
    ? singleLine
    : `${singleLine.slice(0, maxLength - 2)}...`;
});

// https://stackoverflow.com/a/37511463
defineStringProperty("removeAccents", function () {
  return this.normalize("NFD").replace(/[\u0300-\u036f]/gu, "");
});

defineStringProperty("toSearchValue", function () {
  return this.toLowerCase().removeAccents();
});

defineStringProperty("fuzzyMatch", function (search) {
  if (!search) return true;

  if (search.includes(" ")) {
    this.toSearchValue().includes(search.toSearchValue());
  }

  return this.toSearchValue()
    .split(" ")
    .some((word) => word.startsWith(search.toSearchValue()));
});

export {};
