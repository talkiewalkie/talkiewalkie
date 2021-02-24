import React, { useRef, useState } from "react";
import { Field, Form, Formik } from "formik";
import { Navigate, useNavigate } from "react-router-dom";
import { useSwipeable } from "react-swipeable";
import Dropzone from "react-dropzone";
import classNames from "classnames";

import { useUser } from "../contexts";
import { apiPostb, upload } from "../utils";
import { useRecorder } from "../hooks/useRecorder";
import { Spinner } from "../components";

export const Editor = () => {
  const { user } = useUser();
  const navigate = useNavigate();
  const handlers = useSwipeable({
    onSwipedLeft: () => navigate("/editor"),
  });
  const { audio, isRecording, startRecording, stopRecording } = useRecorder();

  const audioPreviewRef = useRef<HTMLAudioElement>(null);
  const recorderRef = useRef<HTMLAudioElement>(null);

  const [duration, setDuration] = useState<number>();
  const [isPlaying, setIsPlaying] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!user) return <Navigate to="/" />;

  return (
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
      onSubmit={async ({ title, description, cover, audio }) => {
        setIsSubmitting(true);
        const uploads = await Promise.all([upload(cover!), upload(audio!)]);
        apiPostb("auth/walk/create", {
          title,
          description,
          coverArtUuid: uploads[0].uuid,
          audioUuid: uploads[1].uuid,
        }).then(() => navigate("/"));
      }}
    >
      {({ values, setFieldValue }) =>
        isSubmitting ? (
          <div className="flex h-full">
            <Spinner inline={false} />
          </div>
        ) : (
          <Form
            className="flex-col space-y-16 p-24 m-24 bg-white rounded-sm shadow-sm"
            {...handlers}
          >
            <Dropzone
              accept="image/*"
              multiple={false}
              maxSize={1e7} // 10mb
              onDrop={(fs) => setFieldValue("cover", fs[0])}
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
            {!values.audio && (
              <>
                <button
                  className="mx-auto h-56 w-56 rounded-full shadow-sm-outlined text-24"
                  onClick={async () => {
                    if (isRecording) {
                      stopRecording();
                      if (!audio) return;

                      const af = new File(
                        [audio],
                        `${new Date().getUTCMilliseconds()}`,
                      );
                      setFieldValue("audio", af);
                      const ac = new AudioContext();
                      const ad = await ac.decodeAudioData(
                        await af.arrayBuffer(),
                      );
                      af && setDuration(ad.duration);
                    } else startRecording();
                  }}
                  type="button"
                >
                  {isRecording ? "⏹" : "🎙"}
                </button>
                <audio ref={recorderRef} />
              </>
            )}
            <Dropzone
              accept="audio/*"
              multiple={false}
              maxSize={1e8} // 100mb
              onDrop={async (fs) => {
                setFieldValue("audio", fs[0]);
                const ac = new AudioContext();
                const ad = await ac.decodeAudioData(await fs[0].arrayBuffer());
                setDuration(ad.duration);
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
                    <div>
                      {Math.floor(duration!! / 60)}:
                      {(Math.floor(duration!!) % 60)
                        .toString()
                        .padStart(2, "0")}
                    </div>
                  ) : (
                    "Upload audio from file"
                  )}
                  <input {...getInputProps()} />
                </div>
              )}
            </Dropzone>
            <input type="submit" />
          </Form>
        )
      }
    </Formik>
  );
};
