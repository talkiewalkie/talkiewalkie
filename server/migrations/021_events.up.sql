CREATE TYPE "event_type" AS ENUM (
    'new_message',
    'deleted_message',
    'changed_picture',
    'joined_conversation',
    'left_conversation',
    'conversation_title_changed'
    );

CREATE TABLE "event"
(
    "id"                   SERIAL PRIMARY KEY,
    "uuid"                 UUID                                           NOT NULL DEFAULT "uuid_generate_v4"(),
    "created_at"           TIMESTAMP WITH TIME ZONE                       NOT NULL DEFAULT now(),
    "type"                 "event_type"                                   NOT NULL,
    "recipient_id"         INT REFERENCES "user" ("id") ON DELETE CASCADE NOT NULL,

    "message_id"           INT                                            REFERENCES "message" ("id") ON DELETE SET NULL,
    CONSTRAINT "message_id_should_be_null_unless_type" CHECK ( "message_id" IS NULL OR
                                                               "type" = 'new_message'::"event_type"
        ),
        
    "conversation_id"      INT                                            REFERENCES "conversation" ("id") ON DELETE SET NULL,
    CONSTRAINT "conversation_id_should_be_null" CHECK ( "conversation_id" IS NULL OR
                                                        "type" = 'joined_conversation'::"event_type" OR
                                                        "type" = 'left_conversation'::"event_type" OR
                                                        "type" = 'conversation_title_changed'::"event_type"
        ),

    "deleted_message_uuid" UUID,
    CONSTRAINT "deleted_message_uuid_should_be_null" CHECK ( "deleted_message_uuid" IS NULL OR "type" = 'deleted_message' )

);

ALTER TABLE "user"
    ADD COLUMN "is_online"      BOOL                     NOT NULL DEFAULT FALSE,
    ADD COLUMN "last_online_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now();
