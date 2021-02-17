import { useEffect, useRef } from "react";

const on = <E extends keyof WindowEventMap>(
  type: E,
  listener: (this: Window, ev: WindowEventMap[E]) => any,
  useCapture?: boolean,
) => {
  window.addEventListener(type, listener, useCapture);
  return () => window.removeEventListener(type, listener, useCapture);
};

const useSyncRef = <T>(value: T) => {
  const ref = useRef<T>(value);

  useEffect(() => {
    ref.current = value;
  }, [value]);

  return ref;
};

export const useOn = <E extends keyof WindowEventMap>(
  type: E,
  listener: (ev: WindowEventMap[E]) => any,
  useCapture?: boolean,
) => {
  const listenerRef = useSyncRef(listener);
  useEffect(() => on(type, listenerRef.current, useCapture), [
    type,
    listenerRef,
    useCapture,
  ]);
};
