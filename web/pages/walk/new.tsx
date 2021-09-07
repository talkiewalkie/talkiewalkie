import {
  useAuthUser,
  withAuthUser,
  withAuthUserTokenSSR,
} from "next-firebase-auth";
import withLayout, { Coords, LocationContext } from "../../components/Layout";
import { FormProvider, useForm, useFormContext } from "react-hook-form";
import { poster } from "../../lib/api";
import ReactMapGL, { Marker } from "react-map-gl";
import { LocationMarkerIcon } from "@heroicons/react/outline";
import { CheckCircleIcon, DotsHorizontalIcon } from "@heroicons/react/solid";
import React, {
  ChangeEvent,
  HTMLProps,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from "react";
import { useContextOrThrow } from "../../lib/useContext";
import Modal from "../../components/Modal";
import { FormFileInput } from "../../components/FormFileInput";
import { useRouter } from "next/router";

const CoverAndWalkPicker = () => {
  const [length, setLength] = useState<string>();

  return (
    <div className="flex w-full text-sm">
      <FormFileInput name="cover" className="w-1/2 max-h-48 cursor-pointer">
        {(files) => (
          <>
            <label className="hidden" htmlFor="cover">
              Picture
            </label>
            {files && files[0] ? (
              <div className="relative">
                <img
                  className="w-full max-h-48 object-cover"
                  src={URL.createObjectURL(files[0])}
                  alt="selected cover picture"
                />
                <div className="absolute top-0 left-0 h-full w-full bg-black opacity-20 hover:opacity-40" />
              </div>
            ) : (
              <div className="bg-gray-300 hover:bg-gray-200 h-48 flex justify-center items-center">
                <p className="text-sm text-gray-500">Upload your cover 📷</p>
              </div>
            )}
          </>
        )}
      </FormFileInput>
      <FormFileInput
        name="walk"
        className="w-1/2 bg-blue-500 hover:bg-blue-400 p-2 text-white cursor-pointer"
        onChange={(e) => {
          if (!e.target.files || !e.target.files[0]) return;

          const reader = new FileReader();
          const ac = new AudioContext();
          reader.onload = (fe) => {
            fe.target?.result &&
              // @ts-ignore
              ac.decodeAudioData(fe.target.result, (buffer) =>
                setLength(
                  buffer.duration > 59
                    ? `${Math.floor(buffer.duration / 60)}min`
                    : `${Math.floor(buffer.duration)}s`
                )
              );
          };
          reader.readAsArrayBuffer(e.target.files[0]);
        }}
        accept="audio/*"
        multiple={false}
      >
        {(files) => (
          <label className="truncate" htmlFor="walk">
            {files && files[0]
              ? `🔉 ${length} - ${files[0].name}`
              : "Pick your track 🔉"}
          </label>
        )}
      </FormFileInput>
    </div>
  );
};

const NewWalk = () => {
  const { position: userPosition } = useContextOrThrow(LocationContext);
  const authUser = useAuthUser();
  const { push } = useRouter();

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
          onSubmit={handleSubmit(({ title, description, cover, walk }) => {
            const fd = new FormData();
            fd.append("cover", cover[0]);
            fd.append("walk", walk[0]);
            fd.append(
              "payload",
              JSON.stringify({ title, description, startPoint: position })
            );
            poster("/walk", fd)
              .then((res) => res.json())
              .then((data) =>
                push({ pathname: "/walk/[uuid]", query: { uuid: data.uuid } })
              );
          })}
        >
          <div className="px-4 flex items-center">
            <div className="h-8 w-8 bg-red-400 rounded-full mr-4" />
            <span className="hover:underline">{authUser.displayName}</span>
            <div className="ml-auto">
              <DotsHorizontalIcon height={16} />
            </div>
          </div>

          <CoverAndWalkPicker />

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
