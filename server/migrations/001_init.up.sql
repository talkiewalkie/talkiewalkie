CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE "user"
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    handle      CITEXT      NOT NULL UNIQUE,
    email       CITEXT      NOT NULL UNIQUE,
    password    TEXT        NOT NULL UNIQUE,
    email_token TEXT
);

CREATE TABLE "walk"
(
    id        SERIAL PRIMARY KEY,
    uuid      UUID UNIQUE                                    NOT NULL DEFAULT uuid_generate_v4(),
    title     TEXT                                           NOT NULL,
    cover_url TEXT,
    author_id INT REFERENCES "user" ("id") ON DELETE CASCADE NOT NULL
);