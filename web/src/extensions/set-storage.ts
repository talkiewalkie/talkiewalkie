declare global {
  interface Window {
    E2E: true | undefined;
    storage: Storage;
  }
  const storage: Storage;
}

window.storage = window.E2E ? sessionStorage : localStorage;

export {};
