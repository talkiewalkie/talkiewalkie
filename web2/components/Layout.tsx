import React, { ReactNode, useEffect, useState } from "react";
import Head from "next/head";
import Link from "next/link";
import { useAuthUser } from "next-firebase-auth";
import firebase from "firebase/app";
import "firebase/auth";
import { MailIcon } from "@heroicons/react/solid";
import { useForm } from "react-hook-form";
import Modal from "./Modal";

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
                .catch((e) => setAuthError(e));
            else
              firebase
                .auth()
                .createUserWithEmailAndPassword(form.email, form.password)
                .catch((e) => setAuthError(e))
                .then(onSuccess);
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

function Layout({ children }: { children: ReactNode }) {
  const user = useAuthUser();
  const [authCardType, setAuthCardType] = useState<"login" | "signup">();
  const [showUserModal, setShowUserModal] = useState(false);

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
              <div className="ml-auto">
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
          <div>{user.email}</div>
          <button
            className="mt-8 text-blue-400 font-medium hover:text-blue-300"
            onClick={user.signOut}
          >
            Sign Out
          </button>
        </Modal>
      </main>
    </div>
  );
}

export default function withLayout(component: () => ReactNode) {
  const Thing = () => <Layout>{component()}</Layout>;

  return Thing;
}
