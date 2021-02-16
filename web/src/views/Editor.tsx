import { Field, Form, Formik } from "formik";
import React from "react";
import { apiPost, apiPostb, upload } from "../utils";

export const Editor = () => {
  return (
    <div>
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
        // validate={({ file }) => ({ file: !!file })}
        onSubmit={async ({ title, description, cover, audio }) => {
          const uploads = await Promise.all([upload(cover!), upload(audio!)]);
          apiPostb("auth/walk/create", {
            title,
            description,
            coverArtUuid: uploads[0].uuid,
            audioUuid: uploads[1].uuid,
          });
        }}
      >
        {({ setFieldValue }) => (
          <Form className="flex-col space-y-16">
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
    </div>
  );
};
