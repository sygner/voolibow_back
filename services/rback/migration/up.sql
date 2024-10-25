DROP TABLE IF EXISTS "roles";
CREATE TABLE "roles" (
    "role_id"       text NOT NULL,
    "description"          text,
    PRIMARY KEY ("role_id")
);

DROP TABLE IF EXISTS "permissions";
CREATE TABLE "permissions" (
    "permission_id" text NOT NULL,
    "description"          text,
    PRIMARY KEY ("permission_id")
);

DROP TABLE IF EXISTS "user_roles";
CREATE TABLE "user_roles" (
    "user_id"       text NOT NULL,
    "role_id"       text NOT NULL REFERENCES "roles" ("role_id") ON DELETE RESTRICT,
    PRIMARY KEY ("user_id", "role_id")      --user is allowed to have more than one role
);

DROP TABLE IF EXISTS "role_permissions";
CREATE TABLE "role_permissions" (
    "role_id"               text NOT NULL REFERENCES "roles" ("role_id") ON DELETE CASCADE,
    "permission_id"         text NOT NULL REFERENCES "permissions" ("permission_id") ON DELETE RESTRICT,
    PRIMARY KEY ("role_id", "permission_id")
);