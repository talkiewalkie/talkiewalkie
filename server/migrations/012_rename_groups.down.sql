ALTER TABLE "conversation"
    RENAME TO "group";
ALTER TABLE "user_conversation"
    RENAME TO "user_group";
ALTER TABLE "user_group"
    RENAME "conversation_id" TO "group_id";
ALTER TABLE "message"
    RENAME "conversation_id" TO "group_id";

