import Link from "next/link";
import { useEffect, useState } from "react";
import { DotsHorizontalIcon } from "@heroicons/react/solid";
import Layout from "../components/Layout";
import useSWR from "swr";

type Walk = {
  uuid: string;
  title: string;
  description: string;
  author: { uuid: string; handle: string };
  coverUrl: string;
};

const WalkCard = ({ walk }: { walk: Walk }) => {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <Link as={`/walk/${walk.uuid}`} href="/walk/[uuid]">
      <div
        key={walk.uuid}
        className="rounded-sm bg-white border py-4 hover:shadow-lg flex flex-col space-y-4"
      >
        <div className="px-4 flex items-center">
          <div className="h-8 w-8 bg-red-400 rounded-full mr-4" />
          <Link as={`/user/${walk.author.uuid}`} href="/user/[uuid]">
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
          <span className="ml-auto text-gray-300">
            {Math.floor(Math.random() ** 2 * 30)}min away
          </span>
        </div>

        <div className="px-4 truncate">
          <span className="text-gray-500">{walk.description}</span>
        </div>
      </div>
    </Link>
  );
};

export default function Home() {
  const { error, data: walks } = useSWR<Walk[]>("http://localhost:8080/walks");
  const loading = !error && !walks;

  return loading ? (
    <div className="h-full mx-auto">loading</div>
  ) : walks ? (
    walks.map((w) => <WalkCard key={w.uuid} walk={w} />)
  ) : null;
}
