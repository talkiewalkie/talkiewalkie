ALTER TABLE "asset"
    ADD COLUMN "bucket"    TEXT,
    ADD COLUMN "blob_name" TEXT,
    ADD CONSTRAINT "bucket_and_blob_name_need_to_be_defined_together" CHECK (("bucket" IS NULL) = ("blob_name" IS NULL));

