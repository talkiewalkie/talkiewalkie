import { withAuthUser, withAuthUserTokenSSR } from "next-firebase-auth";
import withLayout, { Coords, LocationContext } from "../../components/Layout";
import { useDropzone } from "react-dropzone";
import { useForm } from "react-hook-form";
import { poster } from "../../lib/fetcher";
import ReactMapGL, { Marker } from "react-map-gl";
import { LocationMarkerIcon } from "@heroicons/react/outline";
import { XCircleIcon } from "@heroicons/react/solid";
import React, { useState } from "react";
import { useContextOrThrow } from "../../lib/useContext";

const NewWalk = () => {
  const [viewport, setViewport] = useState({
    latitude: 48.853,
    longitude: 2.3499,
    zoom: 10,
  });
  const { position: userPosition } = useContextOrThrow(LocationContext);
  const [position, setPosition] = useState<Coords>(userPosition);

  const { handleSubmit, register } = useForm({});
  const { acceptedFiles, getRootProps, getInputProps } = useDropzone();

  return (
    <form
      className="flex flex-col space-y-4"
      onSubmit={handleSubmit((data) => poster("walk", data))}
    >
      <label htmlFor="picture">Picture</label>
      <div {...getRootProps({ className: "dropzone" })}>
        <input {...register("picture")} {...getInputProps()} />
        <p>Drag 'n' drop some files here, or click to select files</p>
      </div>
      <label htmlFor="walk">Audio</label>
      <input type="file" {...register("walk")} />
      <label>Position</label>
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

      <label htmlFor="title">Title</label>
      <input
        type="text"
        placeholder="Title of your walk..."
        {...register("title")}
      />
      <label htmlFor="description">Description</label>
      <input
        type="textarea"
        placeholder="Description of your walk..."
        {...register("description")}
      />

      <button type="submit">Make it!</button>
    </form>
  );
};

export const getServerSideProps = withAuthUserTokenSSR()();

export default withAuthUser()(withLayout(NewWalk));
