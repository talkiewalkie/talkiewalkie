ALTER TABLE "message"
    -- we directly store the protobuf message sent by the client here since the backend won't perform any db
    -- operation based on the transcript's contents.
    ADD COLUMN "siri_transcript" BYTEA,
    ADD CONSTRAINT "transcript_should_be_null_unless_voice" CHECK ( "type" = 'voice'::MESSAGE_TYPE OR "siri_transcript" IS NULL );
