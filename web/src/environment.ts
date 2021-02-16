export const isDev = window.location.host.includes("localhost");
export const apiUrl = isDev
  ? "http://localhost:8080"
  : "https://api.talkiewalkie.app";
