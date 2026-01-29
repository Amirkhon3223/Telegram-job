-- Add post type enum
CREATE TYPE post_type AS ENUM ('vacancy', 'resume');

-- Add employment type enum for resumes
CREATE TYPE employment_type AS ENUM ('full-time', 'part-time', 'contract', 'freelance');

-- Rename jobs table to posts
ALTER TABLE jobs RENAME TO posts;

-- Rename indexes
ALTER INDEX idx_jobs_status RENAME TO idx_posts_status;
ALTER INDEX idx_jobs_created_at RENAME TO idx_posts_created_at;

-- Add post_type column with default 'vacancy' for existing records
ALTER TABLE posts ADD COLUMN post_type post_type NOT NULL DEFAULT 'vacancy';

-- Make company_id nullable (resumes don't have companies)
ALTER TABLE posts ALTER COLUMN company_id DROP NOT NULL;

-- Add resume-specific fields
ALTER TABLE posts ADD COLUMN experience_years NUMERIC(3,1);  -- e.g. 1.5 years
ALTER TABLE posts ADD COLUMN employment employment_type;
ALTER TABLE posts ADD COLUMN about TEXT;                     -- "О кандидате"
ALTER TABLE posts ADD COLUMN resume_link TEXT;               -- link to CV
ALTER TABLE posts ADD COLUMN contact TEXT;                   -- for resumes: direct contact

-- Add user_id to posts (for resumes that don't have company)
ALTER TABLE posts ADD COLUMN user_id UUID REFERENCES users(id);

-- Update existing posts to have user_id from company
UPDATE posts p SET user_id = c.user_id FROM companies c WHERE p.company_id = c.id;

-- Create index for post_type
CREATE INDEX idx_posts_type ON posts(post_type);
