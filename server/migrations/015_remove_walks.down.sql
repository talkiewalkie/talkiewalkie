CREATE TABLE "walk"
(
    "id"          SERIAL
        CONSTRAINT "walk_pkey"
            PRIMARY KEY,
    "uuid"        UUID                     DEFAULT "uuid_generate_v4"() NOT NULL
        CONSTRAINT "walk_uuid_key"
            UNIQUE,
    "title"       TEXT                                                  NOT NULL,
    "cover_id"    INTEGER
        CONSTRAINT "walk_cover_id_fkey"
            REFERENCES "asset",
    "audio_id"    INTEGER
        CONSTRAINT "walk_audio_id_fkey"
            REFERENCES "asset",
    "author_id"   INTEGER                                               NOT NULL
        CONSTRAINT "walk_author_id_fkey"
            REFERENCES "user"
            ON DELETE CASCADE,
    "created_at"  TIMESTAMP WITH TIME ZONE DEFAULT now()                NOT NULL,
    "start_point" POINT                                                 NOT NULL,
    "description" TEXT,
    CONSTRAINT "cover_and_audio_not_same_file"
        CHECK ((("cover_id" IS NULL) <> ("audio_id" IS NULL)) OR ("cover_id" <> "audio_id"))
);

CREATE TABLE "user_walk"
(
    "id"      SERIAL
        CONSTRAINT "user_walk_pkey"
            PRIMARY KEY,
    "user_id" INTEGER NOT NULL
        CONSTRAINT "user_walk_user_id_fkey"
            REFERENCES "user"
            ON DELETE CASCADE,
    "walk_id" INTEGER NOT NULL
        CONSTRAINT "user_walk_walk_id_fkey"
            REFERENCES "walk"
            ON DELETE CASCADE,
    CONSTRAINT "user_walk_user_id_walk_id_key"
        UNIQUE ("user_id", "walk_id")
);