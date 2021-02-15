import React, { useState } from "react";

import { firebaseApp } from "./firebase";

export type User = {
  email: string;
  handle: string;
  uuid: string;
};

type UserContextType = {
  user?: User;
  login?: (email: string, password: string) => void;
  logout?: () => void;
};

const UserC = React.createContext<UserContextType>({});
export const useUser = () => React.useContext(UserC);

export const UserContext = ({ children }: { children: React.ReactNode }) => {
  const [userAuth, setUserAuth] = useState<User>();

  return (
    <UserC.Provider
      value={{
        user: userAuth,
        login: (email, password) =>
          fetch("http://localhost:8080/unauth/login", {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ email, password }),
          })
            .then((r) => r.json())
            .then((u) => setUserAuth(u))
            .catch((r) => console.log(r)),
        logout: () =>
          firebaseApp
            ?.auth()
            ?.signOut()
            ?.then(() => setUserAuth(undefined)),
      }}
    >
      {children}
    </UserC.Provider>
  );
};
