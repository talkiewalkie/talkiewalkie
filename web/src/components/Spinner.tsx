import React from "react";
import classNames from "classnames";

import styles from "./spinner.module.css";

export const Spinner = ({
  className,
  inline,
  small,
}: {
  className?: string;
  inline?: boolean;
  small?: boolean;
}) => (
  <span
    role="status"
    className={classNames(styles.spinner, className, {
      [styles.center]: !inline,
      [styles.small]: small,
    })}
  />
);
