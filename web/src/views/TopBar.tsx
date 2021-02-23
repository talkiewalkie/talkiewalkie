import React, { useState } from "react";
import { Link, useMatch, useNavigate } from "react-router-dom";
import { Field, Form, Formik } from "formik";

import { useAuth, useUser } from "../contexts";
import { Popover } from "../components/Popover";
import { useTargetState } from "../hooks";

export const TopBar = () => {
  const { login, logout } = useAuth();
  const { user, updateCachedUser } = useUser();
  const navigate = useNavigate();

  const [showLogin, setShowLogin] = useState(false);

  const isCreatingAccount = useMatch("/signin");
  const [target, setTarget] = useTargetState();

  return (
    <div className="flex items-center justify-between relative px-24 py-12 border-b shadow-lg">
      <Link to="/">TalkieWalkie</Link>
      {!isCreatingAccount && (
        <>
          <button onClick={setTarget} className="">
            {user ? user.handle : "login/create"}
          </button>

          {target && (
            <Popover
              target={target}
              position="bottom-right"
              onClose={() => setTarget(undefined)}
              className="absolute text-right top-44 right-0 z-1 p-12 bg-white border rounded-sm shadow-lg text-body"
              style={{ width: user ? 100 : 200 }}
            >
              {user ? (
                <div>
                  <div>{user.email}</div>
                  <button
                    onClick={() => {
                      logout();
                      updateCachedUser(undefined);
                      setTarget(undefined);
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
                    setTarget(undefined);
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
                  <Link
                    className="text-left"
                    to="/create"
                    onClick={() => setTarget(undefined)}
                  >
                    create
                  </Link>
                  <button
                    className="text-left"
                    onClick={() => setShowLogin(true)}
                  >
                    login
                  </button>
                </div>
              )}
            </Popover>
          )}
        </>
      )}
    </div>
  );
};
