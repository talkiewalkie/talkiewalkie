import React, { useCallback, useEffect, useState } from "react";
import classNames from "classnames";
import ReactDOM from "react-dom";

import { useClickOutside, useOn } from "../hooks";

import styles from "./popover.module.css";

type ArrowPosition = "top" | "bottom" | "right" | "left";
export type PopoverPosition =
  | ArrowPosition
  | "top-right"
  | "top-left"
  | "bottom-right"
  | "bottom-left";

type FakeElement = { getBoundingClientRect(): ClientRect }; // Allow RichMentionEditor to use coords in state
export type PopoverProps = {
  target: FakeElement;
  position?: PopoverPosition | PopoverPosition[];
  noArrow?: boolean;
  onClose: (() => void) | undefined;
  allowScrolling?: boolean;
} & Omit<React.HTMLAttributes<HTMLDivElement>, "ref">;

export const Popover = ({
  target,
  position = "bottom",
  noArrow,
  onClose,
  allowScrolling,
  className,
  style,
  ...rest
}: PopoverProps) => {
  const ref = useClickOutside(() => onClose?.());

  const getPosition = useCallback(
    () =>
      computePosition({
        target,
        positions: Array.isArray(position) ? position : [position],
        noArrow,
      }),
    [noArrow, position, target]
  );

  const [{ arrowDirection, translate, ...insets }, setPosition] = useState(
    getPosition
  );
  useEffect(() => {
    setPosition(getPosition());
  }, [getPosition]);

  useOn(
    "scroll",
    () => (allowScrolling ? setPosition(getPosition()) : onClose?.()),
    true
  );

  return ReactDOM.createPortal(
    <div
      ref={ref}
      style={{ ...insets, ...style }}
      className={classNames(className, translate, styles.arrowBox, {
        [styles.bottom]: arrowDirection === "bottom",
        [styles.left]: arrowDirection === "left",
        [styles.right]: arrowDirection === "right",
        [styles.top]: arrowDirection === "top",
      })}
      {...rest}
    />,
    document.getElementById("popover-root")!
  );
};

const computePosition = ({
  target,
  positions,
  noArrow,
}: {
  target: FakeElement;
  positions: PopoverPosition[];
  noArrow?: boolean;
}) => {
  const {
    top,
    bottom,
    left,
    right,
    height,
    width,
  } = target.getBoundingClientRect();

  // Computing the position that optimize the available space in two directions
  // Dump x + y for now
  const spaces = {
    top,
    left,
    right: window.innerWidth - right,
    bottom: window.innerHeight - bottom,
  };
  const position = positions.sortDesc((p) => {
    const [axisA, axisB] = p.split("-") as [
      ArrowPosition,
      "right" | "left" | undefined
    ];
    return (
      spaces[axisA] +
      (axisB === "right"
        ? right
        : axisB === "left"
        ? left
        : // Popover is centered. Computing best between top or top-right
        // would require to compute the actual width of the popover.
        // Instead, applying a small avantage of half the width of the target
        // is a good shortcut to favorise centered position
        // if the target is in the middle of the screen
        axisA === "top" || axisA === "bottom"
        ? Math.min(spaces.left, spaces.right) + 1.5 * width
        : Math.min(spaces.top, spaces.bottom) + 1.5 * height)
    );
  })[0];

  const arrowDirectionMap: Partial<
    { [key in PopoverPosition]: ArrowPosition }
  > = {
    top: "bottom",
    bottom: "top",
    left: "right",
    right: "left",
  };
  const arrowDirection = noArrow ? undefined : arrowDirectionMap[position];
  const arrowMargin = arrowDirection ? 11 : 0; // Should match $popover-arrow-size

  const shouldCenterVertically = position === "left" || position === "right";
  const shouldCenterHorizontally = position === "bottom" || position === "top";

  return {
    arrowDirection,
    translate: shouldCenterVertically
      ? "transform -translate-y-1/2"
      : shouldCenterHorizontally
      ? "transform -translate-x-1/2"
      : undefined,
    top: position.includes("bottom")
      ? bottom + arrowMargin
      : shouldCenterVertically
      ? top + height / 2
      : undefined,
    bottom: position.includes("top")
      ? window.innerHeight - top + arrowMargin
      : undefined,
    left:
      position === "right"
        ? right + arrowMargin
        : position === "bottom-left" || position === "top-left"
        ? left
        : shouldCenterHorizontally
        ? left + width / 2
        : undefined,
    right:
      position === "left"
        ? window.innerWidth - left + arrowMargin
        : position === "bottom-right" || position === "top-right"
        ? window.innerWidth - right
        : undefined,
  };
};
