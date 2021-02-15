export const removeCookie = (name: string) => {
  document.cookie = `${name}=;path=/;max-age=-1`;
};
