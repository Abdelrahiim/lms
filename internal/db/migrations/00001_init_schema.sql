-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create custom types
CREATE TYPE user_role AS ENUM ('student', 'instructor', 'admin');
CREATE TYPE enrollment_status AS ENUM ('active', 'completed', 'suspended', 'dropped');
CREATE TYPE course_level AS ENUM ('beginner', 'intermediate', 'advanced');
CREATE TYPE content_type AS ENUM ('video', 'text', 'pdf', 'audio', 'interactive');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS content_type;
DROP TYPE IF EXISTS course_level;
DROP TYPE IF EXISTS enrollment_status;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd