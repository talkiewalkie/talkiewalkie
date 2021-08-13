import "../styles/globals.css";
import type { AppProps } from "next/app";

import initAuth from "../auth";
import { LocationContextProvider } from "../components/Layout";

initAuth();

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <LocationContextProvider>
      <Component {...pageProps} />
    </LocationContextProvider>
  );
}

export default MyApp;
