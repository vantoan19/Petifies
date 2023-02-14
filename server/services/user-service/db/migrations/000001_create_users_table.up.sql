CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   email VARCHAR (300) UNIQUE NOT NULL,
   password VARCHAR (300) NOT NULL,
   first_name VARCHAR(50) NOT NULL,
   last_name VARCHAR(50) NOT NULL,
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);
