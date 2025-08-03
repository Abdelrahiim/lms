-- +goose Up
-- +goose StatementBegin
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL, -- COMP101
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL, -- URL-friendly version
    description TEXT,
    syllabus TEXT,
    instructor_id UUID NOT NULL REFERENCES users(id),
    category VARCHAR(100),
    sub_category VARCHAR(100),
    level VARCHAR(50), -- beginner, intermediate, advanced
    language VARCHAR(10) DEFAULT 'en',
    thumbnail_url VARCHAR(500),
    intro_video_url VARCHAR(500),
    duration_hours INTEGER, -- Estimated duration
    price DECIMAL(10,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'USD',
    is_free BOOLEAN DEFAULT true,
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMP,
    is_featured BOOLEAN DEFAULT false,
    enrollment_type VARCHAR(50) DEFAULT 'open', -- open, approval, invite
    max_students INTEGER,
    prerequisites TEXT[],
    tags TEXT[],
    learning_outcomes TEXT[],
    requirements TEXT[],
    target_audience TEXT,
    completion_certificate BOOLEAN DEFAULT true,
    allow_discussion BOOLEAN DEFAULT true,
    allow_download BOOLEAN DEFAULT false,
    metadata JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}', -- Course-specific settings
    rating_average DECIMAL(3,2) DEFAULT 0.00,
    rating_count INTEGER DEFAULT 0,
    enrolled_count INTEGER DEFAULT 0,
    completed_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_courses_instructor ON courses(instructor_id);
CREATE INDEX idx_courses_published ON courses(is_published, deleted_at);
CREATE INDEX idx_courses_category ON courses(category, sub_category);
CREATE INDEX idx_courses_slug ON courses(slug);
CREATE INDEX idx_courses_search ON courses USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Course Co-instructors & Teaching Assistants
CREATE TABLE course_staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- co_instructor, teaching_assistant
    permissions JSONB DEFAULT '{}', -- Specific permissions for this course
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(id),
    UNIQUE(course_id, user_id)
);

-- Create indexes for course_staff
CREATE INDEX idx_course_staff_course ON course_staff(course_id);
CREATE INDEX idx_course_staff_user ON course_staff(user_id);

-- Modules
CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL,
    is_published BOOLEAN DEFAULT true,
    unlock_type VARCHAR(50) DEFAULT 'immediate', -- immediate, scheduled, sequential
    unlock_date TIMESTAMP,
    prerequisites UUID[], -- Array of module IDs
    estimated_duration_minutes INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(course_id, order_index)
);

-- Create indexes for modules
CREATE INDEX idx_modules_course ON modules(course_id, order_index);

-- Lessons
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(50) NOT NULL, -- video, text, pdf, audio, interactive
    content JSONB NOT NULL, -- Flexible content storage
    order_index INTEGER NOT NULL,
    duration_minutes INTEGER,
    is_preview BOOLEAN DEFAULT false, -- Free preview
    is_published BOOLEAN DEFAULT true,
    allow_comments BOOLEAN DEFAULT true,
    attachments JSONB DEFAULT '[]', -- Array of {name, url, size, type}
    transcript TEXT, -- For videos/audio
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(module_id, order_index)
);

-- Create indexes for lessons
CREATE INDEX idx_lessons_module ON lessons(module_id, order_index);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS course_staff;
DROP TABLE IF EXISTS courses;
-- +goose StatementEnd