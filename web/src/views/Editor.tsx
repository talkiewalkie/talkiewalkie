import React, { useRef, useState } from "react";
import { Field, FormikProvider, useFormik } from "formik";
import { Navigate, useNavigate } from "react-router-dom";
import { useSwipeable } from "react-swipeable";
import Dropzone from "react-dropzone";
import classNames from "classnames";

import { useUser } from "../contexts";
import { apiPostb, upload } from "../utils";
import { useRecorder } from "../hooks";
import { Spinner } from "../components";

type EditorState = {
  title: string;
  description: string;
  cover: File | undefined;
  audio: File | undefined;
};

export const Editor = () => {
  const audioPreviewRef = useRef<HTMLAudioElement>(null);
  const recorderRef = useRef<HTMLAudioElement>(null);

  const [duration, setDuration] = useState<number>();
  const [isPlaying, setIsPlaying] = useState(false);
  const [dzError, setDzError] = useState<string>();

  const { user } = useUser();
  const navigate = useNavigate();
  const handlers = useSwipeable({
    onSwipedLeft: () => navigate("/editor"),
  });

  const formik = useFormik<EditorState>({
    initialValues: {
      cover: undefined,
      audio: undefined,
      title: "",
      description: "",
    },
    onSubmit: async ({ title, description, cover, audio }) => {
      const uploads = await Promise.all([upload(cover!), upload(audio!)]);
      apiPostb("auth/walk/create", {
        title,
        description,
        coverArtUuid: uploads[0].uuid,
        audioUuid: uploads[1].uuid,
      }).then(() => navigate("/"));
    },
    validate: ({ cover, audio, title }) => {
      if (!cover) return { cover: "No cover is selected" };
      if (!audio) return { audio: "No audio is selected" };
      if (title.trim() === "") return { title: "Your walk needs a title!" };
      return {};
    },
  });

  const { handleSubmit, values, errors, setFieldValue, isSubmitting } = formik;

  const { isRecording, startRecording, stopRecording } = useRecorder((blob) => {
    const file = new File([blob], Date.now().toString());
    setFieldValue("audio", file);
    audioFileDuration(file).then(setDuration);
  });

  const error = dzError ?? errors.cover ?? errors.audio ?? errors.title;

  if (!user) return <Navigate to="/" />;
  if (isSubmitting)
    return (
      <div className="flex h-full">
        <Spinner inline={false} />
      </div>
    );

  return (
    <FormikProvider value={formik}>
      <form
        onSubmit={handleSubmit}
        className="flex-col space-y-16 p-24 m-24 bg-white rounded-sm shadow-sm"
        {...handlers}
      >
        {error && <div className="text-danger font-medium">{error}</div>}
        <Dropzone
          accept="image/*"
          multiple={false}
          maxSize={1e7} // 10mb
          onDrop={(fs) => setFieldValue("cover", fs[0])}
          onDropRejected={(fileRejections) =>
            setDzError(fileRejections[0].errors[0].message)
          }
          onDropAccepted={() => setDzError(undefined)}
        >
          {({ getInputProps, getRootProps }) => (
            <div
              className={classNames(
                "flex-center",
                !!values.cover ? "" : "border border-dashed p-16",
              )}
              {...getRootProps()}
            >
              {!!values.cover ? (
                <img
                  className="rounded shadow-lg"
                  src={URL.createObjectURL(values.cover)}
                  alt="uploaded_image"
                />
              ) : (
                "Upload a cover..."
              )}
              <input {...getInputProps()} />
            </div>
          )}
        </Dropzone>
        <label htmlFor="title">Title</label>
        <Field
          type="text"
          name="title"
          placeholder="title"
          className="border p-8"
        />
        {!!values.audio && !!values.cover && (
          <Field
            type="text"
            name="description"
            placeholder="description"
            className="border p-8"
          />
        )}
        <div className="flex items-center space-x-24">
          {values.audio && (
            <>
              <button
                className="mx-auto h-56 w-56 rounded-full bg-white shadow-sm-outlined hover:bg-light-gray"
                onClick={() => {
                  isPlaying
                    ? audioPreviewRef.current?.pause()
                    : audioPreviewRef.current?.play();
                  setIsPlaying(!isPlaying);
                }}
                type="button"
              >
                {isPlaying ? "⏸" : "▶"}
              </button>
              <audio className="hidden" ref={audioPreviewRef}>
                <source src={URL.createObjectURL(values.audio)} />
              </audio>
            </>
          )}

          <button
            className="mx-auto h-56 w-56 rounded-full shadow-sm-outlined text-24"
            onClick={() => (isRecording ? stopRecording() : startRecording())}
            type="button"
          >
            {isRecording ? "⏹" : "🎙"}
          </button>
        </div>
        <audio ref={recorderRef} />

        <Dropzone
          accept="audio/*"
          multiple={false}
          maxSize={1e8} // 100mb
          onDrop={(fs) => {
            setFieldValue("audio", fs[0]);
            audioFileDuration(fs[0]).then(setDuration);
          }}
        >
          {({ getRootProps, getInputProps }) => (
            <div
              className={classNames(
                "flex-center",
                !!values.audio ? "" : "text-13",
              )}
              {...getRootProps()}
            >
              {!!values.audio ? (
                <div>{formatDuration(duration!!)}</div>
              ) : (
                "Upload audio from file"
              )}
              <input {...getInputProps()} />
            </div>
          )}
        </Dropzone>
        <button
          type="submit"
          className="px-32 py-16 rounded shadow-sm-outlined hover:bg-light-gray"
        >
          Post!
        </button>
      </form>
    </FormikProvider>
  );
};

const formatDuration = (d: number) =>
  `${Math.floor(d / 60)}:${(Math.floor(d) % 60).toString().padStart(2, "0")}`;

const audioFileDuration = (f: File) => {
  const ac = new AudioContext();
  return f
    .arrayBuffer()
    .then((ab) => ac.decodeAudioData(ab))
    .then(({ duration: d }) => d);
};
