import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

import { AuthProvider, ErrorBoundary, UserProvider, useUser } from "./contexts";
// todo: optimisation - split unauth screen and authed app: https://kentcdodds.com/blog/authentication-in-react-applications
import { Home, Signin, TopBar } from "./views";

import "./styles/tailwind.output.css";
import { Editor } from "./views/Editor";
import { SWRConfig } from "swr";

const App = () => (
  <ErrorBoundary>
    <Router>
      <AuthProvider>
        <UserProvider>
          <SWRConfig
            value={{
              revalidateOnFocus: false,
              revalidateOnMount: false,
              revalidateOnReconnect: false,
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
        {user && <Route path="/editor" element={<Editor />} />}
      </Routes>
    </div>
  );
};

export default App;
