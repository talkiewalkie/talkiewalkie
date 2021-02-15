import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

import { UserContext } from "./UserContext";
import { TopBar } from "./TopBar";
import { ErrorBoundary } from "./ErrorBoundary";
import { Home } from "./Home";
import "./tailwind.output.css";
import Signin from "./Signin";

const App = () => (
  <ErrorBoundary>
    <Router>
      <UserContext>
        <div className="flex flex-col">
          <TopBar />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/signin" element={<Signin />} />
          </Routes>
        </div>
      </UserContext>
    </Router>
  </ErrorBoundary>
);

export default App;
