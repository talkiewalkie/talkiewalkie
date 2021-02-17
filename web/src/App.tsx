import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { SWRConfig } from "swr";

import { AuthProvider, ErrorBoundary, UserProvider, useUser } from "./contexts";
// todo: optimisation - split unauth screen and authed app: https://kentcdodds.com/blog/authentication-in-react-applications
import { Editor, Home, Signin, TopBar, Walk } from "./views";

import "./styles/tailwind.output.css";

const App = () => (
  <ErrorBoundary>
    <Router>
      <AuthProvider>
        <UserProvider>
          <SWRConfig
            value={{
              revalidateOnFocus: false,
              revalidateOnMount: true,
              revalidateOnReconnect: true,
              refreshWhenOffline: true,
              refreshWhenHidden: true,
              refreshInterval: 0,
            }}
          >
            <AppView />
          </SWRConfig>
        </UserProvider>
      </AuthProvider>
    </Router>
  </ErrorBoundary>
);

const AppView = () => {
  const { user } = useUser();

  return (
    <div className="flex flex-col">
      <TopBar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/create" element={<Signin />} />
        <Route path="/walk/:uuid" element={<Walk />} />
        <Route path="/editor" element={<Editor />} />
      </Routes>
    </div>
  );
};

export default App;
