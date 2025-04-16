-- Migration up: Make bio and date_of_birth nullable
ALTER TABLE job_seekers
    ALTER COLUMN bio DROP NOT NULL,
    ALTER COLUMN date_of_birth DROP NOT NULL;

-- Set default empty string for bio to avoid NULL issues in existing code
UPDATE job_seekers SET bio = '' WHERE bio IS NULL;
