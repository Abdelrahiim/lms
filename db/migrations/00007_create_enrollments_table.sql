-- +goose Up
-- Enrollments
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'active', -- active, completed, suspended, dropped
    enrollment_type VARCHAR(50), -- self, admin, instructor, gift
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    suspended_at TIMESTAMP,
    suspended_reason TEXT,
    dropped_at TIMESTAMP,
    dropped_reason TEXT,
    progress_percentage DECIMAL(5,2) DEFAULT 0.00,
    grade VARCHAR(10), -- A+, A, B+, etc.
    grade_points DECIMAL(5,2),
    certificate_issued BOOLEAN DEFAULT false,
    certificate_issued_at TIMESTAMP,
    certificate_url VARCHAR(500),
    last_accessed_at TIMESTAMP,
    time_spent_minutes INTEGER DEFAULT 0,
    notes TEXT, -- Instructor notes
    metadata JSONB DEFAULT '{}',
    UNIQUE(user_id, course_id)
);

CREATE INDEX idx_enrollments_user ON enrollments(user_id, status);
CREATE INDEX idx_enrollments_course ON enrollments(course_id, status);
CREATE INDEX idx_enrollments_completed ON enrollments(completed_at);

-- Enrollment Requests (for approval-based courses)
CREATE TABLE enrollment_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending', -- pending, approved, rejected
    reason_for_joining TEXT,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, course_id)
);

CREATE INDEX idx_enrollment_requests_course ON enrollment_requests(course_id, status);
-- Access Codes (for invite-only courses)
CREATE TABLE access_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    max_uses INTEGER,
    used_count INTEGER DEFAULT 0,
    valid_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_access_codes_code ON access_codes(code);
CREATE INDEX idx_access_codes_course ON access_codes(course_id);

-- +goose Down
DROP TABLE IF EXISTS enrollment_requests;
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS access_codes;