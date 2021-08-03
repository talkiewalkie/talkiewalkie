import { createContext } from "react";
import { Audio } from "expo-av";

export type Clip = {
  audio: Audio.Sound;
  band: 1 | 2 | 3;
  offsetMs: number;
};

export type EditorContext = {
  markerMs: number;
  playing: boolean;
  clips: Clip[];
};

export const blankEditorContext = {
  markerMs: 0,
  playing: false,
  clips: [],
};

export const EditorContext = createContext<EditorContext>(blankEditorContext);
