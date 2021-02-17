import React from "react";
import useSWR from "swr/esm";
import { useAuth } from "./AuthContext";
import { apiGet } from "../utils";

export type User = {
  uuid: string;
  handle: string;
  email: string;
  verified: boolean;
};

type UserContextType = {
  user: User | undefined;
  updateCachedUser: (data: User | undefined) => void;
};

const UserContext = React.createContext<UserContextType | null>(null);
export const useUser = () => {
  const userContext = React.useContext(UserContext);
  if (!userContext) throw new Error("null user context!");
  return userContext;
};

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { state } = useAuth();

  // todo: bug - `state` changes are not propagated?
  const { data: user, mutate } = useSWR<User>(
    state === "LOGGED_IN" ? "auth/me" : null,
    apiGet
  );

  return (
    <UserContext.Provider
      value={{
        user,
        updateCachedUser: mutate,
      }}
    >
      {children}
    </UserContext.Provider>
  );
};
