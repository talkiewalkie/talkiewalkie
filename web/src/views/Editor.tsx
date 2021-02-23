import { Field, Form, Formik } from "formik";
import React from "react";
import { apiPostb, upload } from "../utils";
import { useUser } from "../contexts";
import { Navigate, useNavigate } from "react-router-dom";
import { useSwipeable } from "react-swipeable";

export const Editor = () => {
  const { user } = useUser();
  const navigate = useNavigate();
  const handlers = useSwipeable({
    onSwipedLeft: () => navigate("/editor"),
  });

  return user ? (
    <Formik<{
      title: string;
      description: string;
      cover: File | undefined;
      audio: File | undefined;
    }>
      initialValues={{
        cover: undefined,
        audio: undefined,
        title: "",
        description: "",
      }}
      // todo: validate files, verify ahead that they're of the expected type
      // validate={({ file }) => ({ file: !!file })}
      onSubmit={async ({ title, description, cover, audio }) => {
        const uploads = await Promise.all([upload(cover!), upload(audio!)]);
        apiPostb("auth/walk/create", {
          title,
          description,
          coverArtUuid: uploads[0].uuid,
          audioUuid: uploads[1].uuid,
        }).then(() => navigate("/"));
      }}
    >
      {({ setFieldValue }) => (
        <Form
          className="flex-col space-y-16 p-24 m-24 bg-white rounded-sm shadow-sm"
          {...handlers}
        >
          <label htmlFor="title">title</label>
          <Field type="text" name="title" placeholder="title" />
          <label htmlFor="description">description</label>
          <Field type="text" name="description" placeholder="description" />
          <label>cover</label>
          <input
            type="file"
            name="cover"
            onChange={(e) => {
              setFieldValue("cover", e.target.files![0]);
            }}
          />
          <label>audio</label>
          <input
            type="file"
            name="audio"
            onChange={(e) => {
              setFieldValue("audio", e.target.files![0]);
            }}
          />
          <input type="submit" />
        </Form>
      )}
    </Formik>
  ) : (
    <Navigate to="/" />
  );
};
