-- Drop indexes
DROP INDEX IF EXISTS idx_job_seekers_user_id;
DROP INDEX IF EXISTS idx_companies_user_id;
DROP INDEX IF EXISTS idx_verification_tokens_token;
DROP INDEX IF EXISTS idx_verification_tokens_user_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop dependent tables first (respecting foreign key constraints)
DROP TABLE IF EXISTS job_seeker_educations;
DROP TABLE IF EXISTS job_seeker_experiences;
DROP TABLE IF EXISTS job_seeker_links;
DROP TABLE IF EXISTS job_seeker_preferences;
DROP TABLE IF EXISTS job_seeker_projects;
DROP TABLE IF EXISTS job_seeker_skills;
DROP TABLE IF EXISTS job_seekers;

-- Drop company-related tables
-- Remove foreign key constraint first
ALTER TABLE IF EXISTS company_addresses DROP CONSTRAINT IF EXISTS fk_company_id;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS company_addresses;

-- Drop remaining tables
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS users;
