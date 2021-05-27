import React, { useCallback, useContext, useEffect, useState } from "react";
import { Text, View } from "react-native";
import { useApi } from "./ApiContext";
import { API } from "./api";

export const Query = <Q extends keyof API["query"]>({
  path,
  params,
  children,
}: {
  path: Q;
  params?: API["query"][Q]["in"];
  children: (out: API["query"][Q]["out"], refetch: () => void) => JSX.Element;
}) => {
  const client = useApi();

  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<any[] | undefined>();
  const [output, setOutput] = useState<API["query"][Q]["out"]>();

  const fetch = useCallback(() => {
    client
      .get<API["query"][Q]["out"]>(path, { params })
      .then((res) => {
        setOutput(res.data);
        setLoading(false);
        setErrors(undefined);
      })
      .catch((err) => {
        setErrors(err);
        setLoading(false);
        setOutput(undefined);
      });
  }, [client, path, params, setOutput, setLoading, setErrors]);

  useEffect(() => {
    setLoading(true);
    fetch();
  }, []);

  return loading ? (
    <View>
      <Text>Loading</Text>
    </View>
  ) : errors ? (
    <View>
      <Text>error</Text>
    </View>
  ) : output === undefined ? (
    <View>
      <Text>no output for this query {path}</Text>
    </View>
  ) : (
    children(output, fetch)
  );
};
