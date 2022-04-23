ALTER TABLE "user"
    ADD COLUMN "handle" CITEXT UNIQUE;

UPDATE "user"
SET "handle" = concat("display_name", "phone_number");

UPDATE "user"
SET "display_name" = coalesce("display_name", "phone_number");

UPDATE "user"
SET "handle" = 'talkiewalkie-theo'
WHERE "firebase_uid" = 'k6WhmQLnpvUCeKuDdpknVzBUu9r1';

ALTER TABLE "user"
    ALTER COLUMN "handle" SET NOT NULL,
    ALTER COLUMN "display_name" SET NOT NULL,
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
--

DELETE
FROM "user"
WHERE "firebase_uid" = 'YUqVmo08xvXqPZLTYXX7qkvuvGn2';
