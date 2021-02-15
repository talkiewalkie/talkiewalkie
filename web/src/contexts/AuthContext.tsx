import React, { useContext, useState } from "react";
import { removeCookie } from "../utils";

type AuthState = {
  state: "LOGGED_IN" | "LOGGED_OUT";
  login: (email: string, password: string) => void;
  logout: () => void;
};

const AuthContext = React.createContext<AuthState | null>(null);

export const useAuth = () => {
  const auth = useContext(AuthContext);
  if (!auth) throw new Error("no auth here");
  return auth;
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const token = document.cookie
    .split("; ")
    .find((c) => c.startsWith("jwt="))
    ?.replace("jwt=", "");

  const [state, setState] = useState<AuthState["state"]>(
    token ? "LOGGED_IN" : "LOGGED_OUT"
  );

  return (
    <AuthContext.Provider
      value={{
        login: (email, password) =>
          fetch("http://localhost:8080/unauth/login", {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ email, password }),
          }).then(() => {
            setState("LOGGED_IN");
          }),
        logout: () => {
          setState("LOGGED_OUT");
          removeCookie("jwt");
        },
        state,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
