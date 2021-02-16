import React from "react";
import { Field, Form, Formik } from "formik";
import { useNavigate } from "react-router-dom";
import { apiPostb } from "../utils";

export const Signin = () => {
  const navigate = useNavigate();

  return (
    <div>
      <Formik
        initialValues={{ email: "", password: "" }}
        validate={({ password }) =>
          password.trim() === "" ? new Error("empty password") : undefined
        }
        onSubmit={({ email, password }) =>
          apiPostb("unauth/user/create", { email, password }).then(() =>
            navigate("/")
          )
        }
      >
        <Form className="flex-col">
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
          <button type="submit">Submit</button>
        </Form>
      </Formik>
    </div>
  );
};

export default Signin;
