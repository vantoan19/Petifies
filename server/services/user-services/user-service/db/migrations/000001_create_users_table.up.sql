CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   email VARCHAR (300) UNIQUE NOT NULL,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   first_name VARCHAR(50),
   last_name VARCHAR(50),
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);
