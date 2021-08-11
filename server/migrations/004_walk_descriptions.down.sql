ALTER TABLE "walk"
    DROP COLUMN "description";

ALTER TABLE "walk"
    ADD COLUMN "end_point" POINT NOT NULL;