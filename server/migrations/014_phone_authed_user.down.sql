ALTER TABLE "user"
    ADD COLUMN "handle" CITEXT UNIQUE,
    -- max phone number size: https://stackoverflow.com/a/4729239
    DROP "phone_number",

    DROP "onboarding_finished",
    DROP "display_name",
    DROP "locales",
    DROP CONSTRAINT "onboarding_attributes_should_have_same_nullity",
    DROP CONSTRAINT "onboarding_attributes_null_on_unfinished";

UPDATE "user"
SET "handle" = "uuid"::TEXT;

ALTER TABLE "user"
    ALTER COLUMN "handle" SET NOT NULL;
