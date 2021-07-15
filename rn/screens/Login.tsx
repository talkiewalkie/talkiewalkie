import React from "react";
import { Text, View } from "react-native";
import { MainNavigatorProps } from "../navigation/MainNavigator";
import { AuthScreen } from "./AuthScreen";
import AllBandButton from "../components/AllBandButton";

export const Login = ({ navigation }: { navigation: MainNavigatorProps }) => (
  <View style={{ display: "flex", flexDirection: "column" }}>
    <Text>This is Talkiewalkie</Text>
    <AllBandButton
      label="Login"
      onPress={() => navigation.navigate("FEED")}
      style={{ marginVertical: 20 }}
    />
    <AllBandButton
      label="Editor"
      onPress={() =>
        navigation.navigate({ name: "EDITOR", params: { uuid: undefined } })
      }
      style={{ marginVertical: 20 }}
    />
    <AuthScreen navigation={navigation} />
  </View>
);
