export type User = {
  uuid: string;
  handle: string;
  email: string;
};

export type Walk = {
  uuid: string;
  title: string;
  author: User;
  coverUrl: string;
};

export type API = {
  query: {
    "unauth/walks": {
      in: never;
      out: (Omit<Walk, "author"> & { author: Pick<User, "uuid" | "handle"> })[];
    };
    "unauth/walk": {
      out: Omit<Walk, "author"> & { author: Pick<User, "uuid" | "handle"> };
    };
    "auth/me": {
      in: never;
      out: {}[];
    };
  };
  mutation: {
    "auth/like": {
      in: { uuid: string };
      out: null;
    };
  };
};
