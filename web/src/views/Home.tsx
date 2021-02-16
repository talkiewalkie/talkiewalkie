import React from "react";
import useSWR from "swr";

import { apiGet } from "../utils";
import { User } from "../contexts";

type Walk = {
  uuid: string;
  title: string;
  coverUrl: string;
  audioUrl: string;
  author: User;
};

export const Home = () => {
  const { data: walks } = useSWR<Walk[]>("unauth/walks", apiGet);

  return (
    <div className="p-12 bg-light-blue">
      {walks !== undefined ? (
        <div className="flex-col space-y-24">
          {walks.length === 0 ? (
            <div>no walks yet</div>
          ) : (
            walks.map((w) => (
              <div
                key={w.uuid}
                className="p-12 flex flex-fill justify-between items-center bg-white shadow-sm-outlined rounded-sm hover:bg-light-gray"
                style={{
                  height: 72,
                  backgroundImage: `linear-gradient(rgba(0, 0, 255, 0.5), rgba(255, 255, 0, 1)), url(${w.coverUrl})`,
                }}
              >
                <div className="flex-col">
                  <div className="font-medium text-white text-18">
                    {w.title}
                  </div>
                  <div>by {w.author.handle}</div>
                </div>
                <audio src={w.audioUrl} controls />
              </div>
            ))
          )}
        </div>
      ) : (
        <div>loading</div>
      )}
    </div>
  );
};
