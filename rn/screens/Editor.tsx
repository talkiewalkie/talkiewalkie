import React, { useEffect, useRef, useState } from "react";
import { Animated, Modal, Text, TouchableOpacity, View } from "react-native";
import { StackScreenProps } from "@react-navigation/stack";
import { MainNavigatorParamList } from "../navigation/MainNavigator";
import { Audio } from "expo-av";
import { BlurView } from "expo-blur";
import styled from "styled-components/native";
import { blankEditorContext, EditorContext } from "../hooks/useEditor";

export const Editor = ({
  route,
  navigation,
}: StackScreenProps<MainNavigatorParamList, "EDITOR">) => (
  <EditorContext.Provider value={blankEditorContext}>
    <WrappedEditor navigation={navigation} route={route} />
  </EditorContext.Provider>
);

export const WrappedEditor = ({
  route,
  navigation,
}: StackScreenProps<MainNavigatorParamList, "EDITOR">) => {
  const uuid = route.params.uuid;

  const [showMgmtModal, setShowMgmtModal] = useState(false);

  const [recording, setRecording] = useState<Audio.Recording>();
  const [sound, setSound] = useState<Audio.Sound>();
  const [duration, setDuration] = useState<number>();
  const [isPlaying, setIsPlaying] = useState(false);

  const recHeight = useRef(new Animated.Value(0)).current;
  recHeight.addListener((v) => setDuration(v.value));

  useEffect(() => {
    Animated.timing(recHeight, {
      toValue: 10000,
      duration: 30,
      // delay: 0,
      // easing: () => 4,
      useNativeDriver: true,
    }).start();
  }, [recHeight]);

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
      {recording && (
        <View
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
            // height: 100,
            height: (duration ?? 2000) / 100,
          }}
        />
      )}
      {/*// TODO: make this draggable */}
      {/*//  https://www.youtube.com/watch?v=tHWGKdpj1rs*/}
      {/*    https://snack.expo.io/@yoobidev/draggable-component*/}
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
              : "rgb(173,127,127)",
            borderWidth: 2,
            borderColor: "rgba(0,0,0,0.24)",
            height: duration / 50,
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
      <TouchableOpacity
        style={{
          position: "absolute",
          right: 30,
          top: 50,

          borderRadius: 40,
          width: 40,
          height: 40,

          backgroundColor: "rgb(255,255,255)",

          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
        onPress={() => setShowMgmtModal(true)}
      >
        <Text>...</Text>
      </TouchableOpacity>
      <View
        style={{
          marginVertical: 100,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Modal
          visible={showMgmtModal}
          animationType="fade"
          style={{
            margin: 0,
          }}
          onRequestClose={() => setShowMgmtModal(false)}
          onStartShouldSetResponder={() => {
            setShowMgmtModal(false);
            console.log("hey");
            return true;
          }}
          transparent
        >
          <BlurView
            tint="default"
            intensity={80}
            style={{
              flex: 1,
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <View
              style={{
                backgroundColor: "white",
                borderRadius: 10,
                padding: 10,
                width: "90%",
                height: "70%",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",

                shadowColor: "rgb(205,147,147)",
                shadowOffset: {
                  width: 0,
                  height: 4,
                },
                shadowOpacity: 0.28,
                shadowRadius: 4.95,
                elevation: 2,
              }}
            >
              <CloseModalButton onPress={() => setShowMgmtModal(false)} />
              <ModalOption>
                <Text>Save draft</Text>
              </ModalOption>
              <ModalOption>
                <Text>Publish</Text>
              </ModalOption>
            </View>
          </BlurView>
        </Modal>
      </View>
    </View>
  );
};

const ModalOption = styled.TouchableOpacity`
  background-color: rgb(255, 255, 255);
  padding: 10px;
  color: rgb(255, 0, 0);
`;

const CloseModalButton = ({ onPress }: { onPress: () => void }) => (
  <TouchableOpacity
    style={{
      position: "absolute",
      top: 16,
      right: 16,
      borderRadius: 30,
      height: 30,
      width: 30,

      backgroundColor: "rgb(177,177,177)",

      shadowColor: "rgb(205,147,147)",
      shadowOffset: {
        width: 0,
        height: 4,
      },
      shadowOpacity: 0.28,
      shadowRadius: 4.95,
      elevation: 2,

      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    }}
    onPress={onPress}
  >
    <Text style={{ color: "white" }}>x</Text>
  </TouchableOpacity>
);
