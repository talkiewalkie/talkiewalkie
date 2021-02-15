import React, { ButtonHTMLAttributes } from "react";

export const Pill = ({
  label,
  ...buttonProps
}: {
  label: string;
} & React.ButtonHTMLAttributes<HTMLButtonElement>) => (
  <button
    className="bg-light-blue hover:bg-white text-body rounded-sm px-16 py-12"
    {...buttonProps}
  >
    {label}
  </button>
);
