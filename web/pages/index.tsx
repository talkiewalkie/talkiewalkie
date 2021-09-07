import { useState } from "react";
import Link from "next/link";
import useSWR from "swr";
import { withAuthUser, withAuthUserTokenSSR } from "next-firebase-auth";
import { DotsHorizontalIcon } from "@heroicons/react/solid";

import withLayout, { LocationContext } from "../components/Layout";
import { fetcher } from "../lib/api";
import { useContextOrThrow } from "../lib/useContext";

type Walk = {
  uuid: string;
  title: string;
  description: string;
  author: { uuid: string; handle: string };
  coverUrl: string;
  distanceFromPoint: number;
};

const displayDistance = (d: number) =>
  d > 10_000
    ? ">10km"
    : d > 950
    ? `>${Math.round(d / 1_000)}km`
    : d > 100
    ? `>${Math.round(d / 100)}00m`
    : "<100m";

const WalkCard = ({ walk }: { walk: Walk }) => {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <Link as={`/walk/${walk.uuid}`} href="/walk/[uuid]" passHref>
      <div
        key={walk.uuid}
        className="rounded-sm bg-white border py-4 hover:shadow-lg flex flex-col space-y-4"
      >
        <div className="px-4 flex items-center">
          <div className="h-8 w-8 bg-red-400 rounded-full mr-4" />
          <Link
            as={`/user/${walk.author.handle}`}
            href="/user/[handle]"
            passHref
          >
            <a className="hover:underline">{walk.author.handle}</a>
          </Link>
          <div className="ml-auto">
            <DotsHorizontalIcon height={16} />
          </div>
        </div>
        <div
          className="relative"
          onMouseEnter={() => setShowDetails(true)}
          onMouseLeave={() => setShowDetails(false)}
        >
          <img
            src={walk.coverUrl}
            className="object-cover"
            alt={walk.title}
            onTouchStart={() => setShowDetails(true)}
            onTouchEnd={() => setShowDetails(false)}
          />
          {showDetails && (
            <>
              <div className="absolute top-0 left-0 h-full w-full bg-black opacity-20" />
              <div className="absolute top-0 left-0 h-full w-full text-white font-medium flex items-center justify-center">
                <div>Length: {Math.floor(Math.random() ** 2 * 30)}min</div>
              </div>
            </>
          )}
        </div>

        <div className="px-4 flex items-center">
          <span>{walk.title}</span>
          <span className="ml-auto text-gray-300 text-sm">
            {displayDistance(walk.distanceFromPoint)} away
          </span>
        </div>

        <div className="px-4 truncate">
          <span className="text-gray-500">{walk.description}</span>
        </div>
      </div>
    </Link>
  );
};

const Home = () => {
  const { position } = useContextOrThrow(LocationContext);
  const { error, data: walks } = useSWR<Walk[]>(["/walks", position], fetcher);

  // TODO: infinite scrolling
  // https://dev.to/hunterjsbit/react-infinite-scroll-in-few-lines-588f
  return (
    <>
      {error ? (
        <div className="h-full mx-auto">{JSON.stringify(error)}</div>
      ) : walks ? (
        <>
          {walks.map((w) => (
            <WalkCard key={w.uuid} walk={w} />
          ))}
        </>
      ) : (
        <div>loading</div>
      )}
    </>
  );
};

export const getServerSideProps = withAuthUserTokenSSR()();

export default withAuthUser()(withLayout(Home));
