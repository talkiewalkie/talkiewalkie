import React, { useContext, useState } from "react";

type AuthContext = {
  state: "LOGGED_IN" | "LOGGED_OUT";
  login: (email: string, password: string) => void;
  logout: () => void;
};

const AuthContext = React.createContext<AuthContext | null>(null);

export const useAuth = () => {
  const auth = useContext(AuthContext);
  if (!auth) throw new Error("no auth here");
  return auth;
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, setState] = useState<AuthContext["state"]>("LOGGED_OUT");

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
        logout: () => null,
        state,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
