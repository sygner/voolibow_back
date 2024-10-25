DROP TABLE IF EXISTS "profiles";
CREATE TABLE "profiles"(
    user_id SERIAL NOT NULL PRIMARY KEY,
    display_sid VARCHAR(40) NOT NULL,
    username TEXT NOT NULL,
    display_username TEXT NOT NULL,
    avatar TEXT NOT NULL,
    updated_at TEXT[] NOT NULL
);