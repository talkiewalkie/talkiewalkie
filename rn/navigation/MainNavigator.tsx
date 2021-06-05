import React from "react";
import {
  createStackNavigator,
  StackNavigationProp,
} from "@react-navigation/stack";
import { Feed } from "../screens/Feed";
import { Login } from "../screens/Login";
import { Walk } from "../screens/Walk";

export type MainNavigatorParamList = {
  LOGIN: undefined;
  FEED: undefined;
  WALK: { uuid: string };
  EDITOR: { uuid?: string };
};

const Stack = createStackNavigator<MainNavigatorParamList>();
export type MainNavigatorProps = StackNavigationProp<MainNavigatorParamList>;

export const MainNavigator = () => {
  return (
    <Stack.Navigator
      screenOptions={{
        // headerShown: false,
        headerLeft: () => null,
        animationEnabled: true,
        gestureDirection: "horizontal",
        gestureEnabled: true,
      }}
    >
      <Stack.Screen name="LOGIN" component={Login} />
      <Stack.Screen
        name="FEED"
        component={Feed}
        options={{ headerTitle: "TALKIWALKI" }}
      />
      <Stack.Screen name="WALK" component={Walk} />
      <Stack.Screen name="EDITOR" component={Feed} />
    </Stack.Navigator>
  );
};
