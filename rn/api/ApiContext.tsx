import React, { createContext, useContext } from "react";
import axios from "axios";

const AxiosContext = createContext(
  axios.create({ url: "https://dev.talkiewalkie.app/" })
);

export const useApi = () => useContext(AxiosContext);

export const ApiProvider = ({ children }: { children: React.ReactNode }) => (
  <AxiosContext.Provider
    value={axios.create({
      baseURL: "https://dev.talkiewalkie.app",
    })}
  >
    {children}
  </AxiosContext.Provider>
);
