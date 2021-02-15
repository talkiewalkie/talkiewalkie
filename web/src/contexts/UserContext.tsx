import React from "react";
import useSWR from "swr/esm";
import { useAuth } from "./AuthContext";

export type User = {
  email: string;
  handle: string;
  uuid: string;
  verified: boolean;
};

type UserContextType = {
  user: User | undefined;
  updateCachedUser: (data: User | undefined) => void;
  query: (endpoint: string) => void;
};

const UserContext = React.createContext<UserContextType | null>(null);
export const useUser = () => {
  const userContext = React.useContext(UserContext);
  if (!userContext) throw new Error("null user context!");
  return userContext;
};

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { state } = useAuth();
  const { data: user, mutate } = useSWR<User>(
    state === "LOGGED_IN" ? "/auth/me" : null,
    () =>
      fetch("http://localhost:8080/auth/me", {
        credentials: "include",
      }).then((r) => r.json()),
    {
      initialData: undefined,
      shouldRetryOnError: false,
    }
  );

  return (
    <UserContext.Provider
      value={{
        user,
        updateCachedUser: mutate,
        query: (endpoint) => fetch(endpoint, { credentials: "include" }),
      }}
    >
      {children}
    </UserContext.Provider>
  );
};
