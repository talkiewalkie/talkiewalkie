import React, { useEffect, useState } from "react";
import { Text, TouchableOpacity, View } from "react-native";
import { StackScreenProps } from "@react-navigation/stack";
import { MainNavigatorParamList } from "../navigation/MainNavigator";
import { Audio } from "expo-av";

export const Editor = ({
  route,
  navigation,
}: StackScreenProps<MainNavigatorParamList, "EDITOR">) => {
  const [recording, setRecording] = useState<Audio.Recording>();
  const [sound, setSound] = useState<Audio.Sound>();
  const [duration, setDuration] = useState<number>();
  const [isPlaying, setIsPlaying] = useState(false);

  const uuid = route.params.uuid;

  useEffect(() => {
    sound
      ?.getStatusAsync()
      ?.then((s) => s.isLoaded && setDuration(s.durationMillis));
  }, [sound]);

  return (
    <View
      style={{
        display: "flex",
        position: "relative",
        flex: 1,
        flexDirection: "row",
      }}
    >
      <View
        style={{
          borderColor: "rgb(255,255,255)",
          borderRightWidth: 1,
          height: "100%",
          flex: 1,
          backgroundColor: "rgb(213,236,210)",
        }}
      />
      <View
        style={{
          borderColor: "#333",
          height: "100%",
          flex: 1,
          backgroundColor: "rgb(219,204,219)",
        }}
      />
      <View
        style={{
          borderColor: "rgb(255,255,255)",
          borderLeftWidth: 1,
          height: "100%",
          flex: 1,
          backgroundColor: "rgb(202,225,227)",
        }}
      />

      {sound && duration && (
        <TouchableOpacity
          style={{
            position: "absolute",
            left: "50%",
            bottom: 100,
            transform: [{ translateX: -30 }],
            borderTopLeftRadius: 20,
            borderTopRightRadius: 20,
            width: 60,
            backgroundColor: isPlaying
              ? "rgb(255, 255, 255)"
              : "rgb(100, 100, 100)",
            height: duration / 100,
          }}
          onPress={async () => {
            if (isPlaying) {
              sound.setStatusAsync({
                shouldPlay: false,
                seekMillisToleranceAfter: 0,
              });
              setIsPlaying(false);
            } else {
              sound
                .setPositionAsync(0)
                .then(() => sound.setStatusAsync({ shouldPlay: true }));
              setIsPlaying(true);
            }
          }}
        />
      )}
      <TouchableOpacity
        style={{
          position: "absolute",
          bottom: 64,
          justifyContent: "center",
          alignItems: "center",
          alignSelf: "center",
          width: 80,
          height: 80,
          left: "50%",
          transform: [{ translateX: -40 }],
          // transform: "translateX(-50%)",
          marginHorizontal: "auto",

          backgroundColor: "#fff",
          padding: 24,
          flex: 0,
          borderRadius: 100,

          elevation: 8,
          shadowColor: "#000",
          shadowOffset: {
            width: 0,
            height: 9,
          },
          shadowOpacity: 0.48,
          shadowRadius: 11.95,
        }}
        onPress={async () => {
          if (recording) {
            setRecording(undefined);
            await recording.stopAndUnloadAsync();
            const uri = recording.getURI();
            if (!uri) throw new Error("could not get uri from recording");
            const { sound } = await Audio.Sound.createAsync(
              { uri },
              { shouldPlay: false }
            );
            sound.setOnPlaybackStatusUpdate(async (status) => {
              if (status.isLoaded && status.didJustFinish) setIsPlaying(false);
            });
            setSound(sound);
          } else {
            await Audio.requestPermissionsAsync();
            await Audio.setAudioModeAsync({
              allowsRecordingIOS: true,
              playsInSilentModeIOS: true,
            });
            const rcd = new Audio.Recording();
            await rcd.prepareToRecordAsync(
              Audio.RECORDING_OPTIONS_PRESET_HIGH_QUALITY
            );
            await rcd.startAsync();
            setRecording(rcd);
          }
        }}
      >
        <Text>{recording ? "p" : "m"}</Text>
      </TouchableOpacity>
    </View>
  );
};
