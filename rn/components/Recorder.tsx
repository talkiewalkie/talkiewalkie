import React, { useState } from "react";
import { Audio } from "expo-av";
import { Text, TouchableOpacity, View } from "react-native";
import { AVPlaybackStatusToSet } from "expo-av/build/Video";

export const Recorder = () => {
  const [recording, setRecording] = useState<Audio.Recording>();
  const [sound, setSound] = useState<Audio.Sound>();
  const [isPlaying, setIsPlaying] = useState(false);

  return (
    <View>
      {sound && (
        <View style={{ flexDirection: "row", justifyContent: "space-between" }}>
          <TouchableOpacity
            style={{
              backgroundColor: "#00ff00",
              marginRight: 18,
              padding: 18,
              borderRadius: 20,
              shadowColor: "#000",
              shadowOffset: {
                width: 0,
                height: 9,
              },
              shadowOpacity: 0.48,
              shadowRadius: 11.95,

              elevation: 18,
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
          >
            <Text>{isPlaying ? "pause" : "play"}</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={{
              backgroundColor: "#00ff00",
              padding: 18,
              borderRadius: 20,
              shadowColor: "#000",
              shadowOffset: {
                width: 0,
                height: 9,
              },
              shadowOpacity: 0.48,
              shadowRadius: 11.95,

              elevation: 18,
            }}
            onPress={() => setSound(undefined)}
          >
            <Text>trash</Text>
          </TouchableOpacity>
        </View>
      )}
      {!sound && (
        <TouchableOpacity
          style={{
            backgroundColor: "#bebebe",
            padding: 18,
            alignSelf: "center",
            borderRadius: 50,
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
                if ("isPlaying" in status && status.didJustFinish)
                  setIsPlaying(false);
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
          <Text>{recording ? "Stop" : "Start"}</Text>
        </TouchableOpacity>
      )}
    </View>
  );
};
