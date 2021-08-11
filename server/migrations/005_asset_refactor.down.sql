ALTER TABLE "asset"
    DROP COLUMN "bucket",
    DROP COLUMN "blob_name",
    DROP CONSTRAINT "bucket_and_blob_name_need_to_be_defined_together";

