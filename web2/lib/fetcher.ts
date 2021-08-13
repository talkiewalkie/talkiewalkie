const api =
  process.env.NODE_ENV === "development"
    ? "http://localhost:8080"
    : "https://api.talkiewalkie.app";

export const fetcher = async (url: string) => {
  const res = await fetch(`${api}${url}`, { credentials: "include" });
  if (res.status > 299) {
    const err = new Error(`Encountered error: ${res.status}`);
    err.name = await res.text();
    throw err;
  }

  return res.json();
};
