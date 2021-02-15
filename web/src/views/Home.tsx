import React, { useEffect, useState } from "react";
import { Editor } from "./Editor";

type Walk = {
  uuid: string;
  title: string;
  cover_url: string;
};

export const Home = () => {
  const [walks, setWalks] = useState<Walk[]>();
  useEffect(() => {
    fetch("http://localhost:8080/unauth/walks")
      .then((r) => r.json())
      .then((b) => {
        setWalks(b as Walk[]);
      });
  }, []);

  return (
    <div className="p-12 bg-light-blue">
      {!!walks ? (
        <div className="flex-col space-y-24">
          {walks.length === 0 ? (
            <div>no walks yet</div>
          ) : (
            walks.map((w) => (
              <div
                key={w.uuid}
                className="p-12 bg-white shadow-sm-outlined hover:bg-light-gray"
              >
                {w.title}
              </div>
            ))
          )}
        </div>
      ) : (
        <div>"loading"</div>
      )}
      <button
        onClick={() =>
          fetch("http://localhost:8080/auth/user/bloub", {
            credentials: "include",
          })
        }
      >
        make auth call
      </button>
      <Editor />
    </div>
  );
};
