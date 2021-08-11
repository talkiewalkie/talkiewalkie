import Layout from "../../components/Layout";
import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import {
  DotsHorizontalIcon,
  PauseIcon,
  PlayIcon,
} from "@heroicons/react/solid";
import Link from "next/link";

type Walk = {
  uuid: string;
  title: string;
  description: string;
  author: { handle: string; uuid: string };
  coverUrl: string;
  audioUrl: string;
};

export default function Walk() {
  const { query, isReady } = useRouter();

  const playerRef = useRef<HTMLAudioElement>(null);

  const [loading, setLoading] = useState(true);
  const [walk, setWalk] = useState<Walk>();
  const [playing, setPlaying] = useState(false);

  useEffect(() => {
    isReady &&
      fetch(`http://localhost:8080/walk/${query.uuid}`)
        .then((r) => r.json())
        .then((d: Walk) => {
          setWalk(d);
          setLoading(false);
        });
  }, [isReady]);

  return loading ? (
    <div>loading</div>
  ) : (
    walk && (
      <div className="rounded-sm bg-white border py-2 flex flex-col space-y-4">
        <audio className="hidden" ref={playerRef}>
          <source src={walk.audioUrl} />
        </audio>
        <div className="px-4 flex items-center">
          <div className="h-8 w-8 bg-red-400 rounded-full mr-4" />
          <Link as={`/user/${walk.author.uuid}`} href="/user/[uuid]">
            <a className="hover:underline">{walk.author.handle}</a>
          </Link>
          <div className="ml-auto">
            <DotsHorizontalIcon height={16} />
          </div>
        </div>
        <div className="relative">
          <img src={walk.coverUrl} alt="cover" />
          <div className="absolute top-0 left-0 h-full w-full bg-black opacity-20" />
          <div className="absolute top-0 left-0 h-full w-full flex items-center justify-center text-white">
            {playing ? (
              <PauseIcon
                height={48}
                onClick={() => {
                  playerRef.current?.pause();
                  setPlaying(false);
                }}
              />
            ) : (
              <PlayIcon
                height={48}
                onClick={() => {
                  playerRef.current?.play();
                  setPlaying(true);
                }}
              />
            )}
          </div>
        </div>
        <div className="px-4 flex items-center">
          <span>{walk.title}</span>
          <span className="ml-auto text-gray-300">
            {Math.floor(Math.random() ** 2 * 30)}min away
          </span>
        </div>

        <div className="px-4">
          <span className="font-medium mr-2">{walk.author.handle}</span>
          <p className="text-gray-500 whitespace-pre-wrap text-sm">
            {walk.description}
          </p>
          <p className="text-xs uppercase text-gray-200 my-4">2min ago</p>
        </div>
      </div>
    )
  );
}
