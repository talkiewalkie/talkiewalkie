import React, { ReactNode, useEffect, useState } from "react";
import Head from "next/head";
import Link from "next/link";
import { useAuthUser } from "next-firebase-auth";
import StyledFirebaseAuth from "react-firebaseui/StyledFirebaseAuth";
import firebase from "firebase";

export default function Layout({ children }: { children: ReactNode }) {
  const user = useAuthUser();
  const [showLogin, setShowLogin] = useState(false);
  const [userPopup, setUserPopup] = useState(false);

  // Do not SSR FirebaseUI, because it is not supported.
  // https://github.com/firebase/firebaseui-web/issues/213
  const [renderAuth, setRenderAuth] = useState(false);

  useEffect(() => {
    if (typeof window !== "undefined") {
      setRenderAuth(true);
    }
  }, []);

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
              <a className="font-bold transform rotate-12">TalkieWalkie</a>
            </Link>
            <div className="ml-auto">
              {user.id ? (
                <>
                  <button
                    className="h-8 w-8 rounded-full bg-gray-300 border"
                    onClick={() => setUserPopup(!userPopup)}
                  />
                  {userPopup && (
                    <div className="fixed top-0 left-0 z-20 bg-green-200">
                      <div>{user.email}</div>
                      <button
                        className="mt-8 text-blue-400 font-medium hover:text-blue-300"
                        onClick={user.signOut}
                      >
                        Sign Out
                      </button>
                    </div>
                  )}
                </>
              ) : (
                <>
                  <button
                    className="rounded p-2 bg-blue-400 text-white font-medium hover:shadow hover:bg-blue-300"
                    onClick={() => setShowLogin(true)}
                  >
                    Log In
                  </button>
                  <button className="p-2 text-blue-400 font-medium hover:text-blue-300">
                    Sign Up
                  </button>
                </>
              )}
              {showLogin && renderAuth && (
                // TODO: style our own
                <StyledFirebaseAuth
                  uiConfig={{
                    signInFlow: "popup",
                    // Auth providers
                    // https://github.com/firebase/firebaseui-web#configure-oauth-providers
                    signInOptions: [
                      {
                        provider: firebase.auth.EmailAuthProvider.PROVIDER_ID,
                        requireDisplayName: true,
                      },
                    ],
                    signInSuccessUrl: "/",
                    credentialHelper: "none",
                    callbacks: {
                      // https://github.com/firebase/firebaseui-web#signinsuccesswithauthresultauthresult-redirecturl
                      signInSuccessWithAuthResult: () => {
                        // Don't automatically redirect. We handle redirecting based on
                        // auth state in withAuthComponent.js.
                        // setShowLogin(false);
                        return false;
                      },
                    },
                  }}
                  firebaseAuth={firebase.auth()}
                />
              )}
            </div>
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
      </main>
    </div>
  );
}

export const withLayout = (component: () => ReactNode) => () => (
  <Layout>{component()}</Layout>
);
