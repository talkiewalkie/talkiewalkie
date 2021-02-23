import React from "react";
import { Link } from "react-router-dom";
import classNames from "classnames";
import hash from "swr/dist/libs/hash";

export type IWalk = {
  uuid: string;
  title: string;
  coverUrl: string;
  author: { uuid: string; handle: string };
};

export const Walk = ({ walk }: { walk: IWalk }) => {
  return (
    <div className="flex-col space-y-8">
      <div className="flex items-center justify-between px-16">
        <Link
          to={`/u/${walk.author.handle}`}
          className="flex items-center space-x-8"
        >
          <div
            className={classNames(
              "h-32 w-32 rounded-full flex-center text-white",
              walk.author.uuid.charCodeAt(0) % 3 === 0
                ? "bg-danger"
                : "bg-light-gray",
            )}
          >
            {walk.author.handle.slice(0, 1)}
          </div>
          <div>{walk.author.handle}</div>
        </Link>
        <button className="">...</button>
      </div>
      <Link
        key={walk.uuid}
        to={`/walk/${walk.uuid}`}
        className="p-16 flex flex-fill justify-between items-center bg-white shadow-sm-outlined hover:bg-light-gray"
        style={{
          height: 112,
          backgroundImage: `linear-gradient(rgba(0, 0, 0, 0.1), rgba(0, 0, 0, 1)), url(${walk.coverUrl})`,
          backdropFilter: "blur(8px)",
        }}
      >
        <div className="flex-col">
          <div className="font-medium text-white text-18">{walk.title}</div>
          <div className="text-white opacity-50">
            {walk.uuid.slice(0, 5)}, {walk.uuid.slice(10, 12)}min away
          </div>
        </div>
      </Link>
    </div>
  );
};
