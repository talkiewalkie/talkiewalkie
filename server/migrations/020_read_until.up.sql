ALTER TABLE "user_conversation"
    ADD COLUMN "read_until" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now();

ALTER TABLE "user"
    ADD COLUMN "last_connected_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now();
