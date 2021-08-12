import { withAuthUser } from "next-firebase-auth";
import { useRouter } from "next/router";
import Link from "next/link";

import withLayout from "../../components/Layout";
import useSWR from "swr";

type User = {
  handle: string;
  bio: string;
  profile: string;
  walks: { uuid: string; title: string }[];
  likes: number;
};

const User = () => {
  const { query, isReady } = useRouter();

  const { error, data: user } = useSWR<User>(() =>
    isReady ? `http://localhost:8080/user/${query.handle}` : null
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
                <div className="font-medium">{user.walks.length}</div>
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
            {user.walks?.map((w) => (
              <Link
                key={w.uuid}
                as={`/walk/${w.uuid}`}
                href="/walk/[uuid]"
                passHref
              >
                <a className="truncate">{w.title}</a>
              </Link>
            ))}
          </div>
        </div>
      ) : (
        <div>loading</div>
      )}
    </div>
  );
};

export default withAuthUser()(withLayout(User));
