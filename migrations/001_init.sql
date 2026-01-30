-- Enum types
CREATE TYPE job_level AS ENUM ('junior', 'middle', 'senior', 'internship', '');
CREATE TYPE job_type AS ENUM ('remote', 'hybrid', 'onsite');
CREATE TYPE job_category AS ENUM ('web2', 'web3', 'dev', '');
CREATE TYPE job_status AS ENUM ('draft', 'pending', 'approved', 'published', 'rejected', 'archived');
CREATE TYPE user_role AS ENUM ('admin', 'recruiter');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');
CREATE TYPE payment_type AS ENUM ('single', 'package', 'subscription');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT,
    role user_role NOT NULL DEFAULT 'recruiter',
    interface_language VARCHAR(5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    contact TEXT NOT NULL,
    telegram TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    level job_level NOT NULL,
    type job_type NOT NULL,
    category job_category NOT NULL,
    salary_from INT,
    salary_to INT,
    description TEXT NOT NULL,
    apply_link TEXT NOT NULL,
    status job_status NOT NULL DEFAULT 'draft',
    language TEXT NOT NULL DEFAULT 'en',
    channel_message_id INT,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id),
    amount INT NOT NULL,
    currency TEXT NOT NULL,
    type payment_type NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
