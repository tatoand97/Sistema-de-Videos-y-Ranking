-- Remove unique constraint or index from users.first_name
ALTER TABLE users DROP CONSTRAINT IF EXISTS ux_users_first_name;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_first_name_key;
DROP INDEX IF EXISTS idx_users_first_name;
