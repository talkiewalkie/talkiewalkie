declare global {
  interface Window {
    E2E: true | undefined;
    storage: Storage;
  }
  const storage: Storage;
}

// eslint-disable-next-line nabla/no-local-storage
window.storage = window.E2E ? sessionStorage : localStorage;

export {};
