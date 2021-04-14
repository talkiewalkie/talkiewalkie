CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE "user"
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    handle      CITEXT      NOT NULL UNIQUE,
    email       CITEXT      NOT NULL UNIQUE,
    password    BYTEA       NOT NULL,
    email_token TEXT
);


CREATE TABLE "asset"
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID UNIQUE              NOT NULL DEFAULT uuid_generate_v4(),
    file_name   TEXT                     NOT NULL,
    mime_type   TEXT                     NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE "walk"
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID UNIQUE                                    NOT NULL DEFAULT uuid_generate_v4(),
    title      TEXT                                           NOT NULL,
    cover_id   INT REFERENCES "asset" ("id"),
    audio_id   INT REFERENCES "asset" ("id"),
    author_id  INT REFERENCES "user" ("id") ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE                       NOT NULL DEFAULT now(),
    CONSTRAINT cover_and_audio_not_same_file CHECK ( ((cover_id IS NULL) <> (audio_id IS NULL)) OR (cover_id <> audio_id) )
);
