ALTER TABLE "message"
    ALTER COLUMN "group_id" DROP NOT NULL,
    DROP CONSTRAINT "message_group_id_fkey",
    ADD CONSTRAINT "message_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE SET NULL
;