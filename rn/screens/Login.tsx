import React from "react";
import { View, Text, TouchableOpacity } from "react-native";
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
    <AuthScreen navigation={navigation} />
  </View>
);
