-- +goose Up
-- Analytics Events
CREATE TABLE analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES user_sessions(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL, -- page_view, video_play, quiz_start, etc.
    event_category VARCHAR(100), -- navigation, content, assessment, etc.
    event_action VARCHAR(100),
    event_label VARCHAR(255),
    event_value INTEGER,
    page_url VARCHAR(500),
    referrer_url VARCHAR(500),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    properties JSONB DEFAULT '{}', -- Custom properties
    user_agent TEXT,
    ip_address INET,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_analytics_events_user ON analytics_events(user_id, created_at);
CREATE INDEX idx_analytics_events_course ON analytics_events(course_id, created_at);
CREATE INDEX idx_analytics_events_type ON analytics_events(event_type, created_at);
CREATE INDEX idx_analytics_events_created ON analytics_events(created_at);

-- Course Analytics (Aggregated)
CREATE TABLE course_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    enrolled_count INTEGER DEFAULT 0,
    active_students INTEGER DEFAULT 0,
    completed_count INTEGER DEFAULT 0,
    dropped_count INTEGER DEFAULT 0,
    average_progress DECIMAL(5,2) DEFAULT 0.00,
    average_time_spent_minutes INTEGER DEFAULT 0,
    lesson_views INTEGER DEFAULT 0,
    quiz_attempts INTEGER DEFAULT 0,
    forum_posts INTEGER DEFAULT 0,
    average_quiz_score DECIMAL(5,2),
    completion_rate DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(course_id, date)
);

CREATE INDEX idx_course_analytics_course ON course_analytics(course_id, date);

-- System Logs
CREATE TABLE system_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    request_method VARCHAR(10),
    request_path VARCHAR(500),
    request_body TEXT,
    response_status INTEGER,
    response_time_ms INTEGER,
    error_message TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_system_logs_user ON system_logs(user_id, created_at);
CREATE INDEX idx_system_logs_action ON system_logs(action, created_at);
CREATE INDEX idx_system_logs_created ON system_logs(created_at);

-- Audit Logs (For compliance)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- File Uploads
CREATE TABLE file_uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uploaded_by UUID NOT NULL REFERENCES users(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    mime_type VARCHAR(100),
    storage_path VARCHAR(500) NOT NULL,
    storage_provider VARCHAR(50) DEFAULT 'local', -- local, s3, gcs
    url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    virus_scanned BOOLEAN DEFAULT false,
    virus_scan_result VARCHAR(50),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_file_uploads_user ON file_uploads(uploaded_by);
CREATE INDEX idx_file_uploads_type ON file_uploads(file_type);

-- Course Ratings
CREATE TABLE course_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enrollment_id UUID NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review_title VARCHAR(255),
    review_content TEXT,
    is_verified_purchase BOOLEAN DEFAULT true,
    helpful_count INTEGER DEFAULT 0,
    unhelpful_count INTEGER DEFAULT 0,
    instructor_response TEXT,
    instructor_response_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(course_id, user_id)
);

CREATE INDEX idx_course_ratings_course ON course_ratings(course_id, rating);
CREATE INDEX idx_course_ratings_user ON course_ratings(user_id);

-- Certificates
CREATE TABLE certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    certificate_number VARCHAR(100) UNIQUE NOT NULL,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    template_id VARCHAR(100),
    verification_url VARCHAR(500),
    pdf_url VARCHAR(500),
    metadata JSONB DEFAULT '{}', -- Additional certificate data
    revoked BOOLEAN DEFAULT false,
    revoked_at TIMESTAMP,
    revoked_reason TEXT
);

CREATE INDEX idx_certificates_enrollment ON certificates(enrollment_id);
CREATE INDEX idx_certificates_number ON certificates(certificate_number);

-- +goose Down
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS course_ratings;
DROP TABLE IF EXISTS file_uploads;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS system_logs;
DROP TABLE IF EXISTS course_analytics;
DROP TABLE IF EXISTS analytics_events;