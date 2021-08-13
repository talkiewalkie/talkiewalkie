import React, { createContext, ReactNode, useContext, useState } from "react";
import Head from "next/head";
import Link from "next/link";
import { useAuthUser } from "next-firebase-auth";
import firebase from "firebase/app";
import "firebase/auth";
import { MailIcon, XCircleIcon } from "@heroicons/react/solid";
import { useForm } from "react-hook-form";
import ReactMapGL, { Marker } from "react-map-gl";
import useSWR from "swr";
import { LocationMarkerIcon } from "@heroicons/react/outline";

import Modal from "./Modal";
import { fetcher } from "../lib/fetcher";
import { useContextOrThrow } from "../lib/useContext";

type Coords = {
  lng: number;
  lat: number;
};

export const LocationContext = createContext<{
  position: Coords;
  setPosition: (p: Coords) => void;
} | null>(null);

export const LocationContextProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const [position, setPosition] = useState<Coords>({
    // Notre-Dame coordinates :P
    lat: 48.853,
    lng: 2.3499,
  });

  return (
    <LocationContext.Provider value={{ position, setPosition }}>
      {children}
    </LocationContext.Provider>
  );
};

const LocationCard = () => {
  const [showMap, setShowMap] = useState(false);
  const [viewport, setViewport] = useState({
    latitude: 48.853,
    longitude: 2.3499,
    zoom: 10,
  });
  const { position, setPosition } = useContextOrThrow(LocationContext);
  const user = useAuthUser();

  return showMap ? (
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
      <button
        className="absolute top-4 right-4 text-gray-500 hover:text-gray-300"
        onClick={() => setShowMap(false)}
      >
        <XCircleIcon height={24} />
      </button>
    </ReactMapGL>
  ) : (
    <div className="flex items-start">
      <button
        className="w-1/2 p-4 font-medium text-white bg-blue-400 hover:bg-blue-300 hover:bg-shadow"
        onClick={() =>
          navigator.geolocation.getCurrentPosition(
            (p) =>
              setPosition({
                lng: p.coords.longitude,
                lat: p.coords.latitude,
              }),
            (e) => setShowMap(true),
            { enableHighAccuracy: true }
          )
        }
      >
        Geoloc
      </button>
      <button
        className="w-1/2 p-4 font-medium text-blue-400 hover:text-blue-300"
        onClick={() => setShowMap(true)}
        disabled={!user.id}
      >
        Set Pin{!user.id && ` (sign in)`}
      </button>
    </div>
  );
};
const AuthCard = ({
  type,
  onSuccess,
}: {
  type: "login" | "signup";
  onSuccess: () => void;
}) => {
  const [emailForm, setEmailForm] = useState(false);
  const [authError, setAuthError] = useState<any>();
  const { register, handleSubmit } = useForm();

  return (
    <div className="p-8 flex flex-col text-left">
      <button
        className="flex items-center space-x-2 hover:opacity-70"
        onClick={() => {
          const googleAuthProvider = new firebase.auth.GoogleAuthProvider();
          firebase.auth().signInWithRedirect(googleAuthProvider);
        }}
      >
        <img
          src="https://upload.wikimedia.org/wikipedia/commons/5/53/Google_%22G%22_Logo.svg"
          className="h-4"
          alt="google-logo"
        />
        <span>{type === "login" ? "Sign In" : "Sign Up"} with Google</span>
      </button>
      <button
        className="mt-4 flex items-center space-x-2 hover:opacity-70"
        onClick={() => setEmailForm(true)}
      >
        <MailIcon height={16} />
        <span>{type === "login" ? "Sign In" : "Sign Up"} with Email</span>
      </button>
      {emailForm && (
        <form
          className="mt-8 flex flex-col space-y-4"
          onSubmit={handleSubmit((form) => {
            if (type === "login")
              firebase
                .auth()
                .signInWithEmailAndPassword(form.email, form.password)
                .catch((e) => setAuthError(e))
                .then(onSuccess);
            else
              firebase
                .auth()
                .createUserWithEmailAndPassword(form.email, form.password)
                .catch((e) => setAuthError(e))
                .then((creds) => {
                  creds &&
                    creds.user?.sendEmailVerification({
                      url: "http://localhost:3000/",
                    });
                  onSuccess();
                });
          })}
        >
          <div className="flex items-center space-x-2">
            <input
              className="rounded-md border focus:border-blue-400 min-w-0 flex-1 px-4 py-2 caret-blue-400"
              type="email"
              placeholder="email"
              autoFocus
              {...register("email")}
            />
          </div>
          <div className="flex items-center space-x-2">
            <input
              className="rounded-md border focus:border-blue-400 min-w-0 flex-1 px-4 py-2 caret-blue-400"
              type="password"
              placeholder="password"
              {...register("password")}
            />
          </div>
          {authError && (
            <span className="text-yellow-600 font-medium">
              {"message" in authError ? authError.message : "Unexpected error"}
            </span>
          )}
          <button
            type="submit"
            className="ml-auto rounded py-1 px-2 font-medium text-white bg-blue-400 hover:bg-blue-300"
            disabled={!!authError}
          >
            {type === "login" ? "Sign In" : "Sign Up"}
          </button>
        </form>
      )}
    </div>
  );
};

type User = {
  handle: string;
};

function Layout({ children }: { children: ReactNode }) {
  const user = useAuthUser();
  const [authCardType, setAuthCardType] = useState<"login" | "signup">();
  const [showUserModal, setShowUserModal] = useState(false);
  const [posPickerModal, setPosPickerModal] = useState(false);

  const { data: me } = useSWR<User>(
    () => (user.clientInitialized && !!user.id && showUserModal ? "/me" : null),
    fetcher
  );

  return (
    <div>
      <Head>
        <title>TalkieWalkie</title>
        <meta name="description" content="TalkieWalkie webapp" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main
        style={{
          minHeight: "100vh",
          height: "100vh",
        }}
      >
        <div className="fixed z-10 w-screen bg-white py-4 border-b border-b-2 flex justify-center">
          <div
            style={{ maxWidth: "40rem" }}
            className="flex items-center w-full px-2"
          >
            <Link href="/">
              <a className="font-bold transform rotate-6 p-1 border border-2 border-yellow-600 rounded">
                TalkieWalkie
              </a>
            </Link>
            {user.clientInitialized && (
              <div className="ml-auto flex items-center">
                <button
                  className="mr-4 text-red-400 hover:text-red-300 animate-pulse"
                  onClick={() => setPosPickerModal(true)}
                >
                  <LocationMarkerIcon height={24} />
                </button>
                {user.id ? (
                  <button
                    className="h-8 w-8 rounded-full bg-gray-300 border"
                    onClick={() => setShowUserModal(!showUserModal)}
                  />
                ) : (
                  <>
                    <button
                      className="rounded p-2 bg-blue-400 text-white font-medium hover:shadow hover:bg-blue-300"
                      onClick={() => setAuthCardType("login")}
                    >
                      Log In
                    </button>
                    <button
                      className="p-2 text-blue-400 font-medium hover:text-blue-300"
                      onClick={() => setAuthCardType("signup")}
                    >
                      Sign Up
                    </button>
                  </>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="bg-gray-100 min-h-full">
          <div
            style={{ maxWidth: "40rem" }}
            className="mx-auto px-2 py-20 h-full flex flex-col space-y-8"
          >
            {children}
          </div>
        </div>

        <Modal open={!!authCardType} onClose={() => setAuthCardType(undefined)}>
          <AuthCard
            type={authCardType!}
            onSuccess={() => setAuthCardType(undefined)}
          />
        </Modal>

        <Modal open={showUserModal} onClose={() => setShowUserModal(false)}>
          <div className="font-medium">{me && me.handle}</div>
          <div className="text-sm text-gray-400">{user.email}</div>
          <button
            className="mt-8 text-blue-400 font-medium hover:text-blue-300"
            onClick={() => user.signOut().then(() => setShowUserModal(false))}
          >
            Sign Out
          </button>
        </Modal>

        <Modal open={posPickerModal} onClose={() => setPosPickerModal(false)}>
          <LocationCard />
        </Modal>
      </main>
    </div>
  );
}

export default function withLayout(component: () => ReactNode) {
  const Thing = () => <Layout>{component()}</Layout>;

  return Thing;
}
