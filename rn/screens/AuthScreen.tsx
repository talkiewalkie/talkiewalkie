import React, { useEffect } from "react";

import * as WebBrowser from "expo-web-browser";
import * as Google from "expo-auth-session/providers/google";
import AllBandButton from "../components/AllBandButton";
import firebase from "firebase";
import "firebase/auth";
import { MainNavigatorProps } from "../navigation/MainNavigator";

if (firebase.apps.length === 0)
  firebase.initializeApp({
    apiKey: "AIzaSyCJUWsDAQY9HnmcdmP6tb8-5wgiIyIt89w",
    authDomain: "talkiewalkie-305117.firebaseapp.com",
    projectId: "talkiewalkie-305117",
    storageBucket: "talkiewalkie-305117.appspot.com",
    messagingSenderId: "719337728540",
    appId: "1:719337728540:web:bbafdbb493565ec3a78dd9",
    measurementId: "G-ZJ2Y9884YN",
  });

WebBrowser.maybeCompleteAuthSession();

export const AuthScreen = ({
  navigation,
}: {
  navigation: MainNavigatorProps;
}) => {
  const [request, response, promptAsync] = Google.useIdTokenAuthRequest({
    clientId:
      "719337728540-q8evdqnbjb6hpqm9pgajki77j9sdb510.apps.googleusercontent.com",
  });

  useEffect(() => {
    if (response?.type === "success") {
      const { id_token } = response.params;

      const credential = firebase.auth.GoogleAuthProvider.credential(id_token);
      firebase
        .auth()
        .signInWithCredential(credential)
        .then(() => navigation.navigate("FEED"));
    }
  }, [response]);

  return (
    <AllBandButton
      label="ToggleAuth"
      onPress={async () => promptAsync()}
      disabled={!request}
    />
  );
};
