import React from "react";
import useSWR from "swr";
import { apiGet } from "../../utils";
import { useNavigate, useParams } from "react-router-dom";

type IWalk = {
  uuid: string;
  title: string;
  coverUrl: string;
  audioUrl: string;
  author: { handle: string; uuid: string };
};

export const Walk = () => {
  const params = useParams();
  const navigate = useNavigate();

  if (!params.hasOwnProperty("uuid")) navigate("/");

  const { data: walk } = useSWR<IWalk>(`unauth/walk/${params.uuid}`, apiGet);
  if (!walk) return <div>spinning</div>;

  return (
    <div className="m-12 p-12 bg-light-blue rounded-sm">
      <div
        className="rounded-sm p-16"
        style={{ height: 400, backgroundImage: `url(${walk.coverUrl})` }}
      >
        <div className="font-medium text-white text-18">{walk.title}</div>
      </div>
    </div>
  );
};
