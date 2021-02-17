import { BaseSyntheticEvent, useCallback, useRef, useState } from "react";
import { useOn } from "./useOn";

export * from "./useOn";

export const useClickOutside = <Element extends HTMLElement = HTMLDivElement>(
  onClickOutside: () => void,
) => {
  const ref = useRef<Element>(null);
  useOn("click", (event) => {
    const element = ref.current;
    if (!element) return;
    if (
      event
        .composedPath()
        .some((target) => target instanceof Element && element.contains(target))
    ) {
      return;
    }
    onClickOutside();
    event.stopPropagation();
  });
  return ref;
};

export const useTargetState = (): [
  Element | undefined,
  (e: BaseSyntheticEvent<any, Element> | undefined) => void,
] => {
  const [target, setTarget] = useState<Element>();

  return [
    target,
    useCallback((e) => {
      const newTarget = e?.currentTarget;
      setTarget(
        newTarget
          ? (current) => (current === newTarget ? undefined : newTarget)
          : undefined,
      );
    }, []),
  ];
};
