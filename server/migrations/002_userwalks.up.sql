CREATE TABLE "user_walk"
(
    id      serial primary key,
    user_id INT REFERENCES "user" (id) ON DELETE CASCADE NOT NULL,
    walk_id INT REFERENCES "walk" (id) ON DELETE CASCADE NOT NULL,
    UNIQUE (user_id, walk_id)
);
