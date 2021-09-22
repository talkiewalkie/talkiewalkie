ALTER TABLE "message"
    ALTER COLUMN "group_id" SET NOT NULL,
    DROP CONSTRAINT "message_group_id_fkey",
    ADD CONSTRAINT "message_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
;