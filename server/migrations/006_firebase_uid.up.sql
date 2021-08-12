ALTER TABLE "user"
    ADD COLUMN "firebase_uid"    VARCHAR(128),
    ADD COLUMN "profile_picture" INT REFERENCES "asset" ("id"),
    DROP COLUMN "password",
    DROP COLUMN "email_token";