ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS fk_sessions_customers;
DROP TABLE IF EXISTS "sessions";