CREATE TABLE "group"
(
    "id"         SERIAL PRIMARY KEY,
    "uuid"       UUID                     NOT NULL DEFAULT "uuid_generate_v4"(),
    "name"       VARCHAR(64),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE "user_group"
(
    "id"         SERIAL PRIMARY KEY,
    "user_id"    INT REFERENCES "user" ("id")  NOT NULL,
    "group_id"   INT REFERENCES "group" ("id") NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE      NOT NULL DEFAULT now(),
    UNIQUE ("user_id", "group_id")
);

CREATE TABLE "message"
(
    "id"         SERIAL PRIMARY KEY,
    "text"       TEXT                     NOT NULL,
    "author_id"  INT REFERENCES "user" ("id"),
    "group_id"   INT REFERENCES "group" ("id"),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);