import React, { useEffect, useState } from "react";

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
      {walks ? (
        <div className="flex-col space-y-24">
          {walks.map((w) => (
            <div
              key={w.uuid}
              className="p-12 bg-white shadow-sm-outlined hover:bg-light-gray"
            >
              {w.title}
            </div>
          ))}
        </div>
      ) : (
        "loading"
      )}
      <button
        onClick={() =>
          fetch("http://localhost:8080/auth/user/bloub", {
            headers: new Headers({
              Authorization:
                document.cookie
                  .split("; ")
                  .find((row) => row.startsWith("fbToken="))
                  ?.split("=")[1] ?? "",
            }),
          })
        }
      >
        {" "}
        make auth call{" "}
      </button>
    </div>
  );
};
