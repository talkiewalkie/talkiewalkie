import { apiUrl } from "./environment";

export const removeCookie = (name: string) => {
  document.cookie = `${name}=;path=/;max-age=-1`;
};

export const apiGet: <T>(
  endpoint: string,
  init?: Omit<RequestInit, "method">,
) => Promise<T> = (endpoint: string, init?: Omit<RequestInit, "method">) =>
  fetch(`${apiUrl}/${endpoint}`, {
    ...init,
    method: "GET",
    credentials: "include",
  })
    .then((r) => {
      if (!r.ok) throw new Error(r.statusText);
      else return r.json();
    })
    .catch((r) => console.log(r));

export const apiPost: <T>(
  endpoint: string,
  body: BodyInit,
  init?: Omit<RequestInit, "method" | "body">,
) => Promise<T> = (
  endpoint: string,
  body: BodyInit,
  init?: Omit<RequestInit, "method" | "body">,
) =>
  fetch(`${apiUrl}/${endpoint}`, {
    ...init,
    method: "POST",
    credentials: "include",
    body,
  })
    .then((r) => {
      if (!r.ok) throw new Error(r.statusText);
      else return r.json();
    })
    .catch((r) => console.log(r));

export const apiPostb = (
  endpoint: string,
  body: Record<string, any> | null,
  init?: Omit<RequestInit, "method" | "body">,
) => apiPost(endpoint, JSON.stringify(body), init);

export const upload = (f: File) => {
  const data = new FormData();
  data.append("main", new Blob([f], { type: f.type }), f.name);
  return apiPost<{ uuid: string; fileName: string; url: string }>(
    "auth/upload",
    data,
  );
};
