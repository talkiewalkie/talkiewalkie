import React, { useState } from "react";
import { useMatch, useNavigate } from "react-router-dom";
import { Field, Form, Formik } from "formik";

import { useAuth, useUser } from "../contexts";
import { Pill } from "../components";

export const TopBar = () => {
  const { login, logout } = useAuth();
  const { user, updateCachedUser } = useUser();
  const navigate = useNavigate();

  const [showUserSpace, setShowUserSpace] = useState(false);
  const [showLogin, setShowLogin] = useState(false);

  const isCreatingAccount = useMatch("/signin");

  return (
    <div className="flex justify-between relative bg-black text-white p-16">
      <button onClick={() => navigate("/")}>TalkieWalkie</button>
      {!isCreatingAccount && (
        <>
          {user && <Pill onClick={() => navigate("/editor")} label="Editor" />}
          <Pill
            onClick={() => setShowUserSpace(!showUserSpace)}
            label={user ? user.handle : "login/create"}
          />

          {showUserSpace && (
            <div className="absolute top-44 right-16 z-1 p-12 bg-white border rounded-sm shadow-lg text-body">
              {user ? (
                <div>
                  <div style={{ display: "absolute" }}>{user.email}</div>
                  <button
                    onClick={() => {
                      logout();
                      updateCachedUser(undefined);
                      setShowUserSpace(false);
                    }}
                  >
                    logout
                  </button>
                </div>
              ) : showLogin ? (
                <Formik
                  initialValues={{ email: "", password: "" }}
                  onSubmit={({ email, password }) => {
                    login(email, password);
                    setShowLogin(false);
                    setShowUserSpace(false);
                    navigate("/");
                  }}
                >
                  <Form className="relative flex-col items-center">
                    <button
                      className="absolute top-0 right-0"
                      onClick={() => setShowLogin(false)}
                      type="button"
                    >
                      x
                    </button>
                    <label htmlFor="email">Email</label>
                    <Field
                      id="email"
                      name="email"
                      placeholder="john@acme.com"
                      type="email"
                    />
                    <label htmlFor="password">Password</label>
                    <Field
                      id="password"
                      name="password"
                      type="password"
                      placeholder="admin1234"
                    />
                    <input type="submit" />
                  </Form>
                </Formik>
              ) : (
                <div className="flex-col">
                  <a className="text-left" href="/create">
                    create
                  </a>
                  <button
                    className="text-left"
                    onClick={() => setShowLogin(true)}
                  >
                    login
                  </button>
                </div>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
};
