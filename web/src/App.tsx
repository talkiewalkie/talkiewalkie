import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

import { ErrorBoundary, AuthProvider, UserProvider } from "./contexts";
import { Home, Signin, TopBar } from "./views";

import "./styles/tailwind.output.css";

const App = () => (
  <ErrorBoundary>
    <Router>
      <AuthProvider>
        <UserProvider>
          <div className="flex flex-col">
            <TopBar />
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/create" element={<Signin />} />
            </Routes>
          </div>
        </UserProvider>
      </AuthProvider>
    </Router>
  </ErrorBoundary>
);

export default App;
