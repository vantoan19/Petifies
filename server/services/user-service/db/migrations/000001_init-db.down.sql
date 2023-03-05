ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS fk_sessions_customers;
DROP TABLE IF EXISTS "sessions";

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS "outbox_status";
DROP TABLE IF EXISTS user_events;