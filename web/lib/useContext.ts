import React, { Context } from "react";

export const useContextOrThrow = <T>(c: Context<T | null>): T => {
  const ctx = React.useContext(c);
  if (!ctx) throw new Error("using context without provider");

  return ctx;
};
