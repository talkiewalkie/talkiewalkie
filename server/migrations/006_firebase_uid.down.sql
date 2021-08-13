ALTER TABLE "user"
    DROP COLUMN "firebase_uid",
    DROP COLUMN "profile_picture",
    ADD COLUMN "password"    BYTEA,
    ADD COLUMN "email_token" TEXT;