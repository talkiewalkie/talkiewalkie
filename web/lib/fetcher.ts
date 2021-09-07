const api =
  process.env.NODE_ENV === "development"
    ? "http://localhost:8080"
    : "https://api.talkiewalkie.app";

export const fetcher = async (url: string, params?: Record<string, any>) => {
  const fullUrl = new URL(`${api}${url}`);
  params &&
    Object.keys(params).forEach((key) =>
      fullUrl.searchParams.append(key, params[key])
    );

  const res = await fetch(fullUrl.toString(), {
    credentials: "include",
    mode: "cors",
  });
  if (res.status > 299) {
    const err = new Error(`Encountered error: ${res.status}`);
    err.name = await res.text();
    throw err;
  }

  return res.json();
};

export const poster = async (url: string, data: FormData) => {
  const res = await fetch(`${api}${url}`, {
    method: "POST",
    credentials: "include",
    mode: "cors",
    body: data,
  });

  if (res.status > 299) {
    const err = new Error(`Encountered error: ${res.status}`);
    err.name = await res.text();
    throw err;
  }

  return res;
};
