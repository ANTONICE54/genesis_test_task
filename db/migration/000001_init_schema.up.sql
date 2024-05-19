CREATE TABLE "emails" (
"id" BIGSERIAL PRIMARY KEY,
"email" VARCHAR NOT NULL UNIQUE,
"created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);




INSERT INTO emails (email) VALUES ('example111@gmail.com');
INSERT INTO emails (email) VALUES ('example222@gmail.com');
INSERT INTO emails (email) VALUES ('example333@gmail.com');