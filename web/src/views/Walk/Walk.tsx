import React, { useRef, useState } from "react";
import classNames from "classnames";
import { Link, useParams } from "react-router-dom";
import useSWR from "swr";

import { apiGet, apiPostb } from "../../utils";
import { useUser } from "../../contexts";
import { Spinner } from "../../components";

type WalkOutput = {
  uuid: string;
  title: string;
  coverUrl: string;
  audioUrl: string;
  author: { handle: string; uuid: string };
  likeCount: number;
  isLiked?: boolean;
};

export const Walk = () => {
  const { user } = useUser();
  const params = useParams();
  const { data: walk, mutate } = useSWR<WalkOutput>(
    user ? `auth/walk/${params.uuid}` : `unauth/walk/${params.uuid}`,
    apiGet,
  );
  const aref = useRef<HTMLAudioElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);

  if (!walk)
    return (
      <div className="h-full flex">
        <Spinner />
      </div>
    );

  return (
    <div className="flex-col m-24 py-16 bg-white rounded shadow-sm">
      <audio className="hidden" ref={aref}>
        <source src={walk.audioUrl} />
      </audio>
      <div className="flex items-center border-b px-16 pb-16 justify-between ">
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
          <div className="flex-col">
            <div className="text-black">{walk.author.handle}</div>
            <div>{walk.author.uuid.slice(0, 4)}, Paris</div>
          </div>
        </Link>
        <button className="">...</button>
      </div>
      <div className="relative bg-light-gray">
        <img
          className="mx-auto"
          style={{
            height: 160,
          }}
          src={walk.coverUrl}
          alt="cover"
        />
        <button
          className="absolute bottom-16 right-16 h-56 w-56 rounded-full bg-white shadow-sm-outlined hover:bg-light-gray"
          onClick={() => {
            isPlaying ? aref.current?.pause() : aref.current?.play();
            setIsPlaying(!isPlaying);
          }}
        >
          {isPlaying ? "⏸" : "▶"}
        </button>
      </div>
      <div className="mt-16 px-16">
        <div className="flex items-center justify-between pb-16">
          <div className="flex-fill font-medium text-18">{walk.title}</div>
          {user && (
            <div className="flex items-center space-x-8">
              <CTA
                label={`${walk.likeCount > 0 ? `${walk.likeCount} ` : ""}${
                  walk.isLiked ? "x" : "♡"
                }`}
                onClick={() =>
                  apiPostb(`auth/walk/like/${walk.uuid}`, null).then(() =>
                    mutate({ isLiked: true, ...walk }),
                  )
                }
                disabled={!!walk.isLiked}
              />
              <CTA label="💬" disabled />
            </div>
          )}
        </div>
        <div
          className="flex-wrap overflow-y-auto text-13"
          style={{ maxHeight: 200 }}
        >
          The purpose of lorem ipsum is to create a natural looking block of
          text (sentence, paragraph, page, etc.) that doesn't distract from the
          layout. A practice not without controversy, laying out pages with
          meaningless filler text can be very useful when the focus is meant to
          be on design, not content. The passage experienced a surge in
          popularity during the 1960s when Letraset used it on their
          dry-transfer sheets, and again during the 90s as desktop publishers
          bundled the text with their software. Today it's seen all around the
          web; on templates, websites, and stock designs. Use our generator to
          get your own, or read on for the authoritative history of lorem ipsum.
        </div>
      </div>
    </div>
  );
};

const CTA = ({
  label,
  ...btnProps
}: { label: string } & Omit<
  React.ButtonHTMLAttributes<HTMLButtonElement>,
  "className"
>) => (
  <button className="bg-light-gray rounded-full px-16 py-4" {...btnProps}>
    {label}
  </button>
);
