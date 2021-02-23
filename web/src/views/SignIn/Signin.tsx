import React from "react";
import { Field, Form, Formik } from "formik";
import { useNavigate } from "react-router-dom";

import { apiPostb } from "../../utils";

export const Signin = () => {
  const navigate = useNavigate();

  return (
    <Formik
      initialValues={{ email: "", password: "", handle: "" }}
      validate={({ password }) =>
        password.trim() === "" ? new Error("empty password") : undefined
      }
      onSubmit={(values) =>
        apiPostb("unauth/user/create", values).then(() => navigate("/"))
      }
    >
      <Form className="p-24 m-24 flex-col space-y-8 bg-white rounded-sm border">
        <label htmlFor="email">Email</label>
        <Field
          id="email"
          name="email"
          placeholder="john@acme.com"
          type="email"
        />
        <label htmlFor="handle">Handle</label>
        <Field id="handle" name="handle" type="text" placeholder="totoLeBg" />
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
  );
};

export default Signin;
