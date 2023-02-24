CREATE TABLE IF NOT EXISTS "sessions" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "user_id" uuid NOT NULL,
    "refresh_token" VARCHAR NOT NULL,
    "exprires_at" timestamptz NOT NULL,
    "client_ip" VARCHAR NOT NULL,
    "is_disabled" BOOLEAN NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" 
    ADD CONSTRAINT fk_sessions_customers
    FOREIGN KEY ("user_id") 
    REFERENCES "users" ("id") 
    ON DELETE CASCADE
    ON UPDATE CASCADE;