import { Form, Formik } from "formik";
import React from "react";

export const Editor = () => {
  return (
    <div>
      <Formik<{ file: File | undefined }>
        initialValues={{ file: undefined }}
        // validate={({ file }) => ({ file: !!file })}
        onSubmit={({ file }) => {
          const fd = new FormData();
          fd.append(
            "main",
            new Blob([file!], { type: file!.type }),
            file!.name
          );
          fetch("http://localhost:8080/auth/upload", {
            method: "POST",
            credentials: "include",
            body: fd,
          }).catch((r) => console.log(r));
        }}
      >
        {({ setFieldValue }) => (
          <Form>
            <label>file</label>
            <input
              type="file"
              name="file"
              onChange={(e) => {
                console.log(e.target.files);
                setFieldValue("file", e.target.files![0]);
              }}
            />
            <input type="submit" />
          </Form>
        )}
      </Formik>
    </div>
  );
};
