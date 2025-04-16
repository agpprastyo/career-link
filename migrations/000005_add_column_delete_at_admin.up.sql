-- migrations/000X_add_deleted_at_to_admins_up.sql
ALTER TABLE admins
ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;

-- Add index to improve query performance when filtering soft-deleted records
CREATE INDEX idx_admins_deleted_at ON admins(deleted_at);

-- Update existing queries to account for soft deletion
CREATE OR REPLACE FUNCTION get_active_admins()
RETURNS SETOF admins AS $$
BEGIN
  RETURN QUERY SELECT * FROM admins WHERE deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;
