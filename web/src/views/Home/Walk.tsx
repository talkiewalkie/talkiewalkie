import React from "react";
import { Link } from "react-router-dom";

export type IWalk = {
  uuid: string;
  title: string;
  coverUrl: string;
  author: { uuid: string; handle: string };
};

export const Walk = ({ walk }: { walk: IWalk }) => {
  return (
    <Link
      key={walk.uuid}
      to={`/walk/${walk.uuid}`}
      className="p-12 flex flex-fill justify-between items-center bg-white shadow-sm-outlined rounded-sm hover:bg-light-gray"
      style={{
        height: 72,
        backgroundImage: `linear-gradient(rgba(0, 0, 255, 0.5), rgba(255, 255, 0, 1)), url(${walk.coverUrl})`,
      }}
    >
      <div className="flex-col">
        <div className="font-medium text-white text-18">{walk.title}</div>
        <div>by {walk.author.handle}</div>
      </div>
    </Link>
  );
};
