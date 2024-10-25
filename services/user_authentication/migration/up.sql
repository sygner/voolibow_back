DROP TABLE IF EXISTS "users";
CREATE TABLE "users"(
    user_id SERIAL NOT NULL PRIMARY KEY,
    phone_number TEXT NOT NULL,
    user_role TEXT NOT NULL,
    user_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

DROP TABLE IF EXISTS "tokens";
CREATE TABLE "tokens"(
    access_token TEXT NOT NULL PRIMARY KEY,
    refresh_token TEXT NOT NULL,
    user_id SERIAL NOT NULL,
    user_role TEXT NOT NULL,
    session_id SMALLINT NOT NULL,
    token_status TEXT NOT NULL,
    ip TEXT NOT NULL,
    agent TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    access_token_expire_at TIMESTAMPTZ NOT NULL,
    refresh_token_expire_at TIMESTAMPTZ NOT NULL
);


DROP TABLE IF EXISTS "login_histories";
CREATE TABLE "login_histories"(
    id TEXT NOT NULL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    user_role TEXT NOT NULL,
    section TEXT NOT NULL,
    ip TEXT NOT NULL,
    agent TEXT NOT NULL,
    logged_in_at TIMESTAMPTZ NOT NULL
);