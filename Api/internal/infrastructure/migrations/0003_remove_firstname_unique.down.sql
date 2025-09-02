-- Recreate unique index on users.first_name
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_first_name ON users(first_name);
