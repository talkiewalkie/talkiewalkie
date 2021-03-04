import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { SWRConfig } from "swr";

import { AuthProvider, ErrorBoundary, UserProvider } from "./contexts";
// todo: optimisation - split unauth screen and authed app: https://kentcdodds.com/blog/authentication-in-react-applications
import { Editor, Home, Signin, TopBar, User, Walk } from "./views";

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

const AppView = () => (
  <div className="h-screen flex-col">
    <TopBar />
    <div className="flex-fill overflow-y-auto bg-light-blue">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/create" element={<Signin />} />
        <Route path="/walk/:uuid" element={<Walk />} />
        <Route path="/u/:handle" element={<User />} />
        <Route path="/editor" element={<Editor />} />
      </Routes>
    </div>
  </div>
);

export default App;
