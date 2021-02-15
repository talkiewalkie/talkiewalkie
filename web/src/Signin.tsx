import React from "react";
import { Field, Form, Formik } from "formik";
import { useNavigate } from "react-router-dom";

export const Signin = () => {
  const navigate = useNavigate();

  return (
    <div>
      <Formik
        initialValues={{ email: "", password: "" }}
        validate={({ password }) =>
          password.trim() === "" ? new Error("empty password") : undefined
        }
        onSubmit={({ email, password }) => {
          fetch("http://localhost:8080/unauth/signin", {
            method: "POST",
            body: JSON.stringify({ email, password }),
          })
            .then(() => navigate("/"))
            .catch((r) => console.log(r));
        }}
      >
        <Form>
          <label htmlFor="email">Email</label>
          <Field
            id="email"
            name="email"
            placeholder="john@acme.com"
            type="email"
          />

          <label htmlFor="password">Password</label>
          <Field id="password" name="password" placeholder="admin1234" />

          <button type="submit">Submit</button>
        </Form>
      </Formik>
    </div>
  );
};

export default Signin;
