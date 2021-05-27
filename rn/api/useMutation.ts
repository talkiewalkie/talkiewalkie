import { API } from "./api";
import { useApi } from "./ApiContext";

const useMutation = <M extends keyof API["mutation"]>(
  path: M
): ((
  params: API["mutation"][M]["in"]
) => Promise<API["mutation"][M]["out"]>) => {
  const client = useApi();

  return (params) => client.post(path, params);
};

// const ee = () => {
//   const like = useMutation("auth/like");
//
//   const out = like({ uuid: "" });
// };
