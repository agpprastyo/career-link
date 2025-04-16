-- Drop trigger and function
DROP TRIGGER IF EXISTS update_users_timestamp ON users;
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop columns
ALTER TABLE users
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS created_at;
