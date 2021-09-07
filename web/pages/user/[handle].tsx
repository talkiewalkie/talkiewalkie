import { withAuthUser } from "next-firebase-auth";
import { useRouter } from "next/router";
import Link from "next/link";

import withLayout from "../../components/Layout";
import useSWR from "swr";
import { fetcher } from "../../lib/api";
import { useState } from "react";

type User = {
  handle: string;
  bio: string;
  profile: string;
  totalWalks: number;
  walks: { uuid: string; title: string }[];
  likes: number;
};

const UserWalks = ({ offset }: { offset: number }) => {
  const { query, isReady } = useRouter();
  const { data: user } = useSWR<User>(
    () => (isReady ? `/user/${query.handle}?offset=${offset}` : null),
    fetcher
  );

  return (
    <>
      {user?.walks?.map((w) => (
        <Link key={w.uuid} as={`/walk/${w.uuid}`} href="/walk/[uuid]" passHref>
          <a className="truncate">{w.title}</a>
        </Link>
      ))}
    </>
  );
};

const User = () => {
  const { query, isReady } = useRouter();
  const [pages, setPages] = useState<number[]>([0]);

  const { error, data: user } = useSWR<User>(
    () => (isReady ? `/user/${query.handle}` : null),
    fetcher
  );

  return (
    <div>
      {error ? (
        <div>{JSON.stringify(error)}</div>
      ) : user ? (
        <div className="bg-white rounded p-2 border flex-col">
          <div className="flex items-center">
            <div className="font-medium">{user.handle}</div>
            <div className="ml-auto flex items-center space-x-4">
              <div className="flex flex-col text-center">
                <div className="font-medium">{user.totalWalks}</div>
                <div className="text-gray-400 text-sm">walks</div>
              </div>
              <div className="flex flex-col text-center">
                <div className="font-medium">{user.likes}</div>
                <div className="text-gray-400 text-sm">likes</div>
              </div>
            </div>
          </div>
          <div className="mt-2 text-gray-400 text-sm whitespace-pre-wrap">
            {user.bio}
          </div>
          <div className="mt-4 flex flex-col space-y-2">
            {pages.map((offset) => (
              <UserWalks key={offset} offset={offset} />
            ))}
            <button
              className="border px-8 py-1 hover:bg-gray-200"
              onClick={() => setPages([...pages, (pages.length + 1) * 20])}
            >
              more
            </button>
          </div>
        </div>
      ) : (
        <div>loading</div>
      )}
    </div>
  );
};

export default withAuthUser()(withLayout(User));
