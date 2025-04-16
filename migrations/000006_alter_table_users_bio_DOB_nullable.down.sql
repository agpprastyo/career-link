-- Migration down: Restore NOT NULL constraints
-- First ensure no NULL values exist
UPDATE job_seekers SET bio = '' WHERE bio IS NULL;
UPDATE job_seekers SET date_of_birth = '2000-01-01'::date WHERE date_of_birth IS NULL;

-- Then restore constraints
ALTER TABLE job_seekers
    ALTER COLUMN bio SET NOT NULL,
    ALTER COLUMN date_of_birth SET NOT NULL;
