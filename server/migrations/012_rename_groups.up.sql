ALTER TABLE "group"
    RENAME TO "conversation";
ALTER TABLE "user_group"
    RENAME TO "user_conversation";
ALTER TABLE "user_conversation"
    RENAME "group_id" TO "conversation_id";
ALTER TABLE "message"
    RENAME "group_id" TO "conversation_id";
