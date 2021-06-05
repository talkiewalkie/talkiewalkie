import { MainNavigatorProps } from "../navigation/MainNavigator";
import {
  ScrollView,
  StyleSheet,
  View,
  Text,
  TouchableOpacity,
} from "react-native";
import React from "react";
import { ApiProvider } from "../api/ApiContext";
import { Query } from "../api/Query";

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignSelf: "stretch",
    paddingHorizontal: 20,
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
  walkCard: {
    backgroundColor: "grey",
    color: "white",
    marginVertical: 5,
    alignSelf: "stretch",
    borderRadius: 4,
    paddingVertical: 10,
    paddingHorizontal: 20,
  },
});

export const Feed = ({ navigation }: { navigation: MainNavigatorProps }) => (
  <View style={styles.container}>
    <ScrollView style={{ flexDirection: "column" }}>
      <ApiProvider>
        <Query path="unauth/walks">
          {(data) => (
            <>
              {data.map((w, idx) => (
                <TouchableOpacity
                  key={w.uuid}
                  style={{
                    ...styles.walkCard,
                    ...(idx === 0 ? { marginTop: 20 } : {}),
                  }}
                  onPress={() => navigation.navigate("WALK", { uuid: w.uuid })}
                >
                  <Text style={{ color: "white" }}>{w.title}</Text>
                  <Text style={{ color: "white" }}>{w.uuid}</Text>
                  <Text
                    style={{
                      color: "white",
                      textAlign: "right",
                      fontStyle: "italic",
                      fontSize: 10,
                    }}
                  >
                    by {w.author.handle}
                  </Text>
                </TouchableOpacity>
              ))}
            </>
          )}
        </Query>
      </ApiProvider>
    </ScrollView>
  </View>
);
