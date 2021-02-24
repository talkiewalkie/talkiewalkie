import React, { useEffect, useState } from "react";

import { useAuth } from "./AuthContext";
import { apiGet } from "../utils";
import { Spinner } from "../components";

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

  const [user, setUser] = useState<User>();
  const [loading, setLoading] = useState(() => state === "LOGGED_IN");

  useEffect(() => {
    if (state === "LOGGED_IN") {
      setLoading(true);
      apiGet<User>("auth/me").then((u) => {
        setUser(u);
        setLoading(false);
      });
    } else {
      setUser(undefined);
      setLoading(false);
    }
  }, [state]);

  return loading ? (
    <div className="h-screen flex">
      <Spinner className="" />
    </div>
  ) : (
    <UserContext.Provider
      value={{
        user,
        updateCachedUser: () => null,
      }}
    >
      {children}
    </UserContext.Provider>
  );
};
