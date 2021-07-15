import React from "react";
import {
  createStackNavigator,
  StackNavigationProp,
} from "@react-navigation/stack";
import { Feed } from "../screens/Feed";
import { Login } from "../screens/Login";
import { Walk } from "../screens/Walk";
import { Editor } from "../screens/Editor";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";

export type MainNavigatorParamList = {
  LOGIN: undefined;
  FEED: undefined;
  WALK: { uuid: string };
  EDITOR: { uuid?: string };
};

const Tab = createBottomTabNavigator<MainNavigatorParamList>();
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
      <Stack.Screen
        name="EDITOR"
        component={Editor}
        options={{ headerShown: false }}
      />
      <Stack.Screen name="WALK" component={Walk} />
    </Stack.Navigator>
  );
};
