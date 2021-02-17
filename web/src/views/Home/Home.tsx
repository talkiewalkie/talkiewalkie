import React from "react";
import useSWR from "swr";

import { apiGet } from "../../utils";
import { IWalk, Walk } from "./Walk";

export const Home = () => {
  const { data: walks } = useSWR<IWalk[]>("unauth/walks", apiGet);

  return (
    <div className="p-12 bg-light-blue">
      {walks !== undefined ? (
        <div className="flex-col space-y-24">
          {walks.length === 0 ? (
            <div>no walks yet</div>
          ) : (
            walks.map((w) => <Walk key={w.uuid} walk={w} />)
          )}
        </div>
      ) : (
        <div>loading</div>
      )}
    </div>
  );
};
