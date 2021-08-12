import { withAuthUser } from "next-firebase-auth";
import { useRouter } from "next/router";

import withLayout from "../../components/Layout";

const User = () => {
  const { query, isReady } = useRouter();
  return <div>{isReady && query.handle}</div>;
};

export default withAuthUser()(withLayout(User));
