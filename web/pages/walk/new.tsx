import {
  useAuthUser,
  withAuthUser,
  withAuthUserTokenSSR,
} from "next-firebase-auth";
import withLayout, { Coords, LocationContext } from "../../components/Layout";
import { DropzoneOptions, useDropzone } from "react-dropzone";
import {
  useForm,
  UseFormRegister,
  Controller,
  useFormContext,
  FormProvider,
} from "react-hook-form";
import { poster } from "../../lib/fetcher";
import ReactMapGL, { Marker } from "react-map-gl";
import { LocationMarkerIcon } from "@heroicons/react/outline";
import { CheckCircleIcon, DotsHorizontalIcon } from "@heroicons/react/solid";
import React, { useCallback, useContext, useEffect, useState } from "react";
import { useContextOrThrow } from "../../lib/useContext";
import Modal from "../../components/Modal";

const WalkPicker = () => {
  const [length, setLength] = useState<string>();
  const { acceptedFiles, getRootProps, getInputProps, fileRejections } =
    useDropzone({
      accept: "audio/*",
      multiple: false,
      noDragEventsBubbling: true,
      onDrop: (files) => {
        const reader = new FileReader();
        const ac = new AudioContext();
        reader.onload = (e) => {
          e.target?.result &&
            // @ts-ignore
            ac.decodeAudioData(e.target.result, (buffer) =>
              setLength(
                buffer.duration > 59
                  ? `${Math.floor(buffer.duration / 60)}min`
                  : `${Math.floor(buffer.duration)}s`
              )
            );
        };
        reader.readAsArrayBuffer(files[0]);
      },
    });

  return (
    <div
      {...getRootProps({
        className:
          "dropzone absolute bottom-8 right-4 transform -rotate-12 bg-blue-500 rounded shadow p-2 text-white text-sm",
      })}
    >
      <input {...getInputProps({ name: "walk" })} />
      {fileRejections[0] ? (
        <p className="text-red-500 font-medium">
          {fileRejections[0].errors[0].message}
        </p>
      ) : acceptedFiles[0] ? (
        <p className="truncate">
          🔉 {length} - {acceptedFiles[0].name}
        </p>
      ) : (
        <p>Pick your track 🔉</p>
      )}
    </div>
  );
};
const CoverAndWalkPicker = ({
  name,
  label,
}: {
  name: string;
  label: string;
}) => {
  // Marrying react-hook-form and dropzone according to
  // https://github.com/react-hook-form/react-hook-form/discussions/2146#discussioncomment-579355
  const { register, unregister, setValue, watch } = useFormContext();
  const file = watch("cover");
  const onDrop = useCallback(
    (files) => setValue("cover", files, { shouldValidate: true }),
    [setValue]
  );
  const { acceptedFiles, getRootProps, getInputProps, fileRejections } =
    useDropzone({
      accept: "image/*",
      multiple: false,
      onDrop,
    });
  useEffect(() => {
    register("cover");
    return () => unregister("cover");
  }, [register, unregister]);

  return (
    <div
      {...getRootProps({
        className: "dropzone max-h-48 cursor-pointer relative",
      })}
    >
      <label className="hidden" htmlFor="picture">
        Picture
      </label>
      <input {...getInputProps({ name })} />
      {fileRejections[0] ? (
        <div className="bg-gray-300 hover:bg-gray-200 h-48 flex justify-center items-center">
          <p className="text-sm text-red-500 font-medium">
            {fileRejections[0].errors.map((e) => e.message).join("\n")}
          </p>
        </div>
      ) : acceptedFiles[0] ? (
        <div className="relative">
          <img
            className="w-full max-h-48 object-cover"
            src={URL.createObjectURL(acceptedFiles[0])}
            alt="selected cover picture"
          />
          <div className="absolute top-0 left-0 h-full w-full bg-black opacity-20 hover:opacity-40" />
        </div>
      ) : (
        <div className="bg-gray-300 hover:bg-gray-200 h-48 flex justify-center items-center">
          <p className="text-sm text-gray-500">{label}</p>
        </div>
      )}
      <WalkPicker />
    </div>
  );
};

const NewWalk = () => {
  const { position: userPosition } = useContextOrThrow(LocationContext);
  const authUser = useAuthUser();

  const [viewport, setViewport] = useState({
    latitude: 48.853,
    longitude: 2.3499,
    zoom: 10,
  });
  const [position, setPosition] = useState<Coords>(userPosition);
  const [showMap, setShowMap] = useState(false);

  const { handleSubmit, register, ...formProps } = useForm({});

  return (
    <div>
      <FormProvider
        handleSubmit={handleSubmit}
        register={register}
        {...formProps}
      >
        <form
          className="rounded-sm bg-white border py-2 flex flex-col space-y-4 shadow-2xl transform scale-75 rotate-3"
          onSubmit={handleSubmit((data) => poster("/walk", data))}
        >
          <div className="px-4 flex items-center">
            <div className="h-8 w-8 bg-red-400 rounded-full mr-4" />
            <span className="hover:underline">{authUser.displayName}</span>
            <div className="ml-auto">
              <DotsHorizontalIcon height={16} />
            </div>
          </div>

          <CoverAndWalkPicker name="cover" label="Pick your cover photo 📷" />

          <div className="flex items-center px-2">
            <label className="text-sm font-medium mr-2" htmlFor="title">
              Title
            </label>
            <input
              type="text"
              className="focus:ring-indigo-500 focus:border-indigo-500 block flex-grow mr-8 px-4 py-2 shadow-sm sm:text-sm border-gray-300 rounded-md"
              placeholder="Title of your walk..."
              {...register("title")}
            />
            <button
              className="ml-auto"
              type="button"
              onClick={() => setShowMap(true)}
            >
              <LocationMarkerIcon height={24} />
            </button>
          </div>

          <div className="p-2 w-full">
            <label className="text-sm font-medium" htmlFor="description">
              Description
            </label>
            <textarea
              className="focus:ring-indigo-500 focus:border-indigo-500 block w-full p-4 shadow-sm sm:text-sm border-gray-300 rounded-md"
              placeholder="Description of your walk..."
              rows={5}
              {...register("description")}
            />
          </div>

          <button type="submit">Make it!</button>
        </form>
      </FormProvider>

      <Modal open={showMap} onClose={() => setShowMap(false)}>
        <div className="flex flex-col space-y-2">
          <ReactMapGL
            mapboxApiAccessToken="pk.eyJ1IjoidGhlby1tLXR3IiwiYSI6ImNrc2FzaTJ5ZzAwMDYycG81bHNvNnU5dmkifQ.VcoUFDpOkF8vJ5TlnpMnJQ"
            {...viewport}
            className="relative"
            width="100%"
            height="300px"
            onViewportChange={(viewport: any) => setViewport(viewport)}
            onClick={(p) => setPosition({ lng: p.lngLat[0], lat: p.lngLat[1] })}
          >
            {position && (
              <Marker longitude={position.lng} latitude={position.lat}>
                <LocationMarkerIcon className="text-blue-400" height={32} />
              </Marker>
            )}
          </ReactMapGL>
          <button
            className="flex justify-center px-16 py-2 rounded bg-blue-500 hover:bg-blue-400 text-white"
            onClick={() => setShowMap(false)}
          >
            <CheckCircleIcon height={24} />
          </button>
        </div>
      </Modal>
    </div>
  );
};

export const getServerSideProps = withAuthUserTokenSSR()();

export default withAuthUser()(withLayout(NewWalk));
