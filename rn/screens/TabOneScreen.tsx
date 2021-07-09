import * as React from "react";
import { Button, StyleSheet, TouchableOpacity } from "react-native";

import EditScreenInfo from "../components/EditScreenInfo";
import { Text, View } from "../components/Themed";
import { ApiProvider } from "../api/ApiContext";
import { Query } from "../api/Query";
import { StackNavigationProp } from "@react-navigation/stack";
import { RootStackParamList } from "../types";
import { BottomTabScreenProps } from "@react-navigation/bottom-tabs";

export const TabOneScreen = () => {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Tab One</Text>
      <Button
        title="Press to call server"
        onPress={() => fetch("http://localhost:8080/unauth/walks")}
      />
      <TouchableOpacity
        onPress={() => fetch("http://localhost:8080/unauth/walks")}
        style={{ padding: 30, backgroundColor: "red" }}
      >
        <Text>Heeey</Text>
      </TouchableOpacity>
      <View
        style={styles.separator}
        lightColor="#eee"
        darkColor="rgba(255,255,255,0.1)"
      />
      <ApiProvider>
        <Query path="walks">
          {(walks, refetch) => (
            <View>
              {walks.map((w) => (
                <View key={w.uuid}>
                  <Text>{w.title}</Text>
                  <TouchableOpacity onPress={refetch}>
                    <Text>Refetch</Text>
                  </TouchableOpacity>
                  {/*<TouchableOpacity onPress={() => navigation.navigate("")}>*/}
                  {/*  <Text>Change screen</Text>*/}
                  {/*</TouchableOpacity>*/}
                </View>
              )) ?? <Text>no walks yet</Text>}
            </View>
          )}
        </Query>
      </ApiProvider>
      <EditScreenInfo path="/screens/TabOneScreen.tsx" />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
  },
  title: {
    fontSize: 20,
    fontWeight: "bold",
  },
  separator: {
    marginVertical: 30,
    height: 1,
    width: "80%",
  },
});
