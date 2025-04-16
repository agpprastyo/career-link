-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'company', 'job_seeker')),
    avatar VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Verification tokens table
CREATE TABLE IF NOT EXISTS verification_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    used_at TIMESTAMP WITH TIME ZONE
);

-- Admins table
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('super', 'admin', 'viewer')),
    last_login TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Company addresses table (created before companies to avoid circular reference)
CREATE TABLE IF NOT EXISTS company_addresses (
    id UUID PRIMARY KEY,
    company_id UUID, -- Will be updated with foreign key after companies table is created
    address_line_1 VARCHAR(255) NOT NULL,
    address_line_2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    zip_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Companies table
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    industry VARCHAR(100),
    website VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    logo_url VARCHAR(255),
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive')),
    size VARCHAR(20) NOT NULL CHECK (size IN ('1-5', '5-10', '10-25', '25-50', '50-100', '100-500', '500-1000', '1000+')),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    address_id UUID NOT NULL REFERENCES company_addresses(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add foreign key constraint to company_addresses
ALTER TABLE company_addresses
    ADD CONSTRAINT fk_company_id
    FOREIGN KEY (company_id)
    REFERENCES companies(id)
    ON DELETE CASCADE;

-- Job seekers table
CREATE TABLE IF NOT EXISTS job_seekers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    bio TEXT NOT NULL,
    profile_picture_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Job seeker skills table
CREATE TABLE IF NOT EXISTS job_seeker_skills (
    id UUID PRIMARY KEY,
    job_seeker_id UUID NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    skill_name VARCHAR(100) NOT NULL,
    proficiency_level VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Job seeker projects table
CREATE TABLE IF NOT EXISTS job_seeker_projects (
    id UUID PRIMARY KEY,
    job_seeker_id UUID NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    link VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Job seeker preferences table
CREATE TABLE IF NOT EXISTS job_seeker_preferences (
    id UUID PRIMARY KEY,
    job_seeker_id UUID NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    preferred_job_type VARCHAR(100) NOT NULL,
    preferred_industry VARCHAR(100) NOT NULL,
    preferred_location VARCHAR(100) NOT NULL,
    salary_expectation INTEGER,
    remote_preference BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Job seeker links table
CREATE TABLE IF NOT EXISTS job_seeker_links (
    id UUID PRIMARY KEY,
    job_seeker_id UUID NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    link_type VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Job seeker experiences table
CREATE TABLE IF NOT EXISTS job_seeker_experiences (
    id UUID PRIMARY KEY,
    job_seeker_id UUID NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    job_title VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Job seeker educations table
CREATE TABLE IF NOT EXISTS job_seeker_educations (
    id UUID PRIMARY KEY,
    job_seeker_id UUID NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    institution_name VARCHAR(255) NOT NULL,
    degree VARCHAR(255) NOT NULL,
    field_of_study VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    grade VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance


CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_job_seekers_user_id ON job_seekers(user_id);
