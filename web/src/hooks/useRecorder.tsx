import { useCallback, useEffect, useState } from "react";

export const useRecorder = () => {
  const [audio, setAudio] = useState<Blob>();
  const [isRecording, setIsRecording] = useState(false);
  const [recorder, setRecorder] = useState<MediaRecorder>();

  const handleData = useCallback((e: BlobEvent) => setAudio(e.data), []);
  useEffect(() => {
    // Lazily obtain recorder first time we're recording.
    if (!recorder) {
      isRecording && requestRecorder().then(setRecorder, console.error);
      return;
    }

    // Manage recorder state.
    if (isRecording) recorder.start();
    else recorder.stop();

    // Obtain the audio when ready.
    recorder.addEventListener("dataavailable", handleData);
    return () => recorder.removeEventListener("dataavailable", handleData);
  }, [isRecording, recorder, handleData]);
  console.log("recostate", recorder?.state);

  return {
    audio,
    isRecording,
    startRecording: () => setIsRecording(true),
    stopRecording: () => setIsRecording(false),
  };
};

const requestRecorder = async () => {
  const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
  return new MediaRecorder(stream);
};
