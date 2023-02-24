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
