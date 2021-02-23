import React from "react";
import useSWR from "swr";
import { useNavigate } from "react-router-dom";
import { useSwipeable } from "react-swipeable";

import { apiGet } from "../../utils";
import { IWalk, Walk } from "./Walk";
import { useUser } from "../../contexts";

export const Home = () => {
  const { data: walks } = useSWR<IWalk[]>("unauth/walks", apiGet);
  const navigate = useNavigate();
  const { user } = useUser();
  const handlers = useSwipeable({
    onSwipedRight: (eventData) => user && navigate("/editor"),
  });

  return (
    <div className="py-12" {...handlers}>
      {walks !== undefined ? (
        walks.length === 0 ? (
          <div>no walks yet</div>
        ) : (
          <div className="flex-col space-y-24 pb-44 mb-44">
            {walks.map((w) => (
              <Walk key={w.uuid} walk={w} />
            ))}
          </div>
        )
      ) : (
        "loading"
      )}
    </div>
  );
};
