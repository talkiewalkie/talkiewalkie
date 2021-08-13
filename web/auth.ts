// ./initAuth.js
import { init } from "next-firebase-auth";

const initAuth = () => {
  init({
    authPageURL: "/auth",
    appPageURL: "/",
    loginAPIEndpoint: "/api/login", // required
    logoutAPIEndpoint: "/api/logout", // required
    firebaseAdminInitConfig: {
      credential: {
        // @ts-ignore
        projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
        // @ts-ignore
        clientEmail: process.env.FIREBASE_CLIENT_EMAIL,
        // @ts-ignore
        privateKey: process.env.FIREBASE_PRIVATE_KEY,
      },
      // @ts-ignore
      databaseURL: process.env.FIREBASE_DATABASE_URL,
    },
    firebaseClientInitConfig: {
      // @ts-ignore
      apiKey: process.env.NEXT_PUBLIC_FIREBASE_PUBLIC_API_KEY,
      authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
      databaseURL: process.env.NEXT_PUBLIC_FIREBASE_DATABASE_URL,
      projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
    },

    cookies: {
      name: "TalkieWalkie", // required
      httpOnly: true,
      domain: "talkiewalkie.app",
      maxAge: 12 * 60 * 60 * 24 * 1000, // twelve days
      overwrite: true,
      path: "/",
      sameSite: "strict",
      secure: process.env.NODE_ENV !== "development", // set this to false in local (non-HTTPS) development
      signed: true,
      keys: [
        // @ts-ignore
        process.env.COOKIE_SECRET_CURRENT,
        // @ts-ignore
        process.env.COOKIE_SECRET_PREVIOUS,
      ],
    },
  });
};

export default initAuth;
