CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   email VARCHAR UNIQUE NOT NULL,
   password VARCHAR NOT NULL,
   first_name VARCHAR NOT NULL,
   last_name VARCHAR NOT NULL,
   is_activated BOOLEAN NOT NULL DEFAULT false,
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);

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

DROP TYPE IF EXISTS "outbox_status";
CREATE TYPE "outbox_state" AS ENUM ('STARTED', 'COMPLETED');
CREATE TABLE IF NOT EXISTS "user_events" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "payload" JSONB NOT NULL,
    "outbox_state" outbox_state NOT NULL,
    "locked_by" uuid NULL,
    "locked_at" timestamptz NULL,
    "error" VARCHAR NULL,
    "completed_at" timestamptz NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);