import { useRouter } from "next/router";
import React, { useRef, useState } from "react";
import {
  DotsHorizontalIcon,
  PauseIcon,
  PlayIcon,
} from "@heroicons/react/solid";
import { LightningBoltIcon } from "@heroicons/react/outline";
import Link from "next/link";
import useSWR from "swr";
import { withAuthUser, withAuthUserTokenSSR } from "next-firebase-auth";

import withLayout, { LocationContext } from "../../components/Layout";
import { fetcher } from "../../lib/api";
import ReactMapGL, { Marker, WebMercatorViewport } from "react-map-gl";
import { LocationMarkerIcon } from "@heroicons/react/outline";
import { useContextOrThrow } from "../../lib/useContext";

type Walk = {
  uuid: string;
  title: string;
  description: string;
  author: { handle: string; uuid: string };
  coverUrl: string;
  audioUrl: string;
  startPoint: { lat: number; lng: number };
};

const WalkMap = ({ walk }: { walk: Walk }) => {
  const { position, setPosition } = useContextOrThrow(LocationContext);
  const { latitude, longitude, zoom } = new WebMercatorViewport({
    width: 800,
    height: 600,
  }).fitBounds(
    [
      [
        Math.min(position.lng, walk.startPoint.lng),
        Math.min(position.lat, walk.startPoint.lat),
      ],
      [
        Math.max(position.lng, walk.startPoint.lng),
        Math.max(position.lat, walk.startPoint.lat),
      ],
    ],
    { padding: 200 }
  );
  const [viewport, setViewport] = useState({
    latitude,
    longitude,
    zoom: Math.min(16, zoom),
  });

  return (
    <ReactMapGL
      mapboxApiAccessToken="pk.eyJ1IjoidGhlby1tLXR3IiwiYSI6ImNrc2FzaTJ5ZzAwMDYycG81bHNvNnU5dmkifQ.VcoUFDpOkF8vJ5TlnpMnJQ"
      {...viewport}
      className="px-4"
      width="100%"
      height="300px"
      onViewportChange={(viewport: any) => setViewport(viewport)}
      onClick={(p) => setPosition({ lng: p.lngLat[0], lat: p.lngLat[1] })}
    >
      {position && (
        <Marker
          longitude={position.lng}
          latitude={position.lat}
          offsetLeft={-16}
          offsetTop={-32}
        >
          <LocationMarkerIcon className="text-blue-400" height={32} />
        </Marker>
      )}
      <Marker
        longitude={walk.startPoint.lng}
        latitude={walk.startPoint.lat}
        offsetLeft={-16}
        offsetTop={-32}
      >
        <LightningBoltIcon className="text-yellow-600" height={32} />
      </Marker>
    </ReactMapGL>
  );
};

const Walk = () => {
  const { query, isReady } = useRouter();

  const playerRef = useRef<HTMLAudioElement>(null);

  const { error, data: walk } = useSWR<Walk>(
    () => (isReady ? `/walk/${query.uuid}` : null),
    fetcher
  );
  const loading = !walk && !error;
  const [playing, setPlaying] = useState(false);

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
          <Link as={`/user/${walk.author.handle}`} href="/user/[handle]">
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
        <WalkMap walk={walk} />
      </div>
    )
  );
};

export const getServerSideProps = withAuthUserTokenSSR()();

export default withAuthUser()(withLayout(Walk));
