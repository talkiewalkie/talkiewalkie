import React, { createContext, useContext } from "react";
import axios from "axios";

const AxiosContext = createContext(
  axios.create({ url: "https://localhost:8080/" })
);

export const useApi = () => useContext(AxiosContext);

export const ApiProvider = ({ children }: { children: React.ReactNode }) => (
  <AxiosContext.Provider
    value={axios.create({
      baseURL: "http://localhost:8080",
    })}
  >
    {children}
  </AxiosContext.Provider>
);
