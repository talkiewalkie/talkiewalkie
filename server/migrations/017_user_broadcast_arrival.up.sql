ALTER TABLE "user"
    ADD COLUMN "broadcast_arrival" BOOL NOT NULL DEFAULT TRUE;
ALTER TABLE "user"
    RENAME COLUMN "bio" TO "status";
