-- migrations/000X_add_deleted_at_to_admins_down.sql
-- Drop function first to avoid dependency issues
DROP FUNCTION IF EXISTS get_active_admins();

-- Drop index
DROP INDEX IF EXISTS idx_admins_deleted_at;

-- Remove column
ALTER TABLE admins
DROP COLUMN IF EXISTS deleted_at;
