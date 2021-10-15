ALTER TABLE "user"
    DROP "handle",
    -- max phone number size: https://stackoverflow.com/a/4729239
    ADD COLUMN "phone_number"        VARCHAR(31) NOT NULL UNIQUE,

    ADD COLUMN "onboarding_finished" BOOL        NOT NULL DEFAULT FALSE,
    ADD COLUMN "display_name"        VARCHAR(64),
    ADD COLUMN "locales"             VARCHAR(5)[],
    ADD CONSTRAINT "onboarding_attributes_should_have_same_nullity" CHECK (("display_name" IS NULL) = ("locales" IS NULL)),
    ADD CONSTRAINT "onboarding_attributes_null_on_unfinished" CHECK (
            ("onboarding_finished" AND "display_name" IS NOT NULL) OR
            (NOT "onboarding_finished" AND "display_name" IS NULL));

