import React from "react";
import { View, Text } from "react-native";
import {
  MainNavigatorParamList,
  MainNavigatorProps,
} from "../navigation/MainNavigator";
import { ApiProvider } from "../api/ApiContext";
import { Query } from "../api/Query";
import { StackScreenProps } from "@react-navigation/stack";

export const Walk = ({
  route,
  navigation,
}: StackScreenProps<MainNavigatorParamList, "WALK">) => {
  return (
    <View>
      <Text>{route.params.uuid}</Text>
      <ApiProvider>
        <Query path={`unauth/walk/${route.params.uuid}`}>
          {(data) => (
            <View>
              <Text>{data.title}</Text>
            </View>
          )}
        </Query>
      </ApiProvider>
    </View>
  );
};
