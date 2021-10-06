CREATE TYPE "message_type" AS ENUM ('text', 'voice', 'image', 'video');

ALTER TABLE "message"
    ADD COLUMN "uuid"         UUID NOT NULL DEFAULT "uuid_generate_v4"(),
    ADD COLUMN "type"         "message_type",
    ADD COLUMN "raw_audio_id" INT REFERENCES "asset" ("id"),
    ALTER COLUMN "text" DROP NOT NULL
;

UPDATE "message"
SET "type" = 'text'::"message_type";

ALTER TABLE "message"
    ALTER COLUMN "type" SET NOT NULL,
    ADD CONSTRAINT "text_message_should_have_non_null_text" CHECK ( "type" <> 'text' OR "text" IS NOT NULL ),
    ADD CONSTRAINT "voice_message_should_have_non_null_raw_audio" CHECK ( "type" <> 'voice' OR "raw_audio_id" IS NOT NULL )
;