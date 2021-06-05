import React from "react";
import { View, Text, TouchableOpacity } from "react-native";
import { MainNavigatorProps } from "../navigation/MainNavigator";

export const Login = ({ navigation }: { navigation: MainNavigatorProps }) => (
  <View style={{ display: "flex" }}>
    <Text>This is Talkiewalkie</Text>
    <TouchableOpacity onPress={() => navigation.navigate("FEED")}>
      <Text>Login</Text>
    </TouchableOpacity>
  </View>
);
