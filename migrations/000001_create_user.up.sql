-- Create enum type for user roles
CREATE TYPE user_role AS ENUM ('admin', 'company', 'job_seeker');

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    avatar VARCHAR(255),
    is_active BOOLEAN DEFAULT FALSE

);

-- Create indexes for frequently queried fields
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
