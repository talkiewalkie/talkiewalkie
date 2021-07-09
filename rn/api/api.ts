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
    walks: {
      in: null;
      out: (Omit<Walk, "author"> & { author: Pick<User, "uuid" | "handle"> })[];
    };
    "unauth/walk": {
      in: { uuid: string };
      out: Omit<Walk, "author"> & { author: Pick<User, "uuid" | "handle"> };
    };
    "auth/me": {
      in: null;
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
