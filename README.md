

# LMS Complete Task Breakdown & Database Design

## Table of Contents
1. [Complete Database Schema](#complete-database-schema)
2. [Authentication & Session Tasks](#authentication--session-tasks)
3. [User Management Tasks](#user-management-tasks)
4. [Course Management Tasks](#course-management-tasks)
5. [Module & Lesson Tasks](#module--lesson-tasks)
6. [Assessment System Tasks](#assessment-system-tasks)
7. [Enrollment Tasks](#enrollment-tasks)
8. [Forum System Tasks](#forum-system-tasks)
9. [Analytics Tasks](#analytics-tasks)
10. [Administration Tasks](#administration-tasks)


## Enhanced Project Structure

```
lms/
├── cmd/
│   ├── api/
│   │   └── main.go                    # API server entry point
│   ├── worker/
│   │   └── main.go                    # Background worker entry point
│   └── migrate/
│       └── main.go                    # Migration runner
│
├── internal/
│   ├── config/
│   │   ├── config.go                  # Configuration struct
│   │   ├── loader.go                  # Environment loader
│   │   └── validator.go               # Config validation
│   │
│   ├── server/
│   │   ├── server.go                  # HTTP server setup
│   │   ├── routes.go                  # Route registration
│   │   └── middleware.go              # Global middleware
│   │
│   ├── database/
│   │   ├── connection.go              # Database connection pool
│   │   ├── transaction.go             # Transaction helpers
│   │   └── sqlc/
│   │       ├── db.go                  # SQLC generated
│   │       ├── models.go              # SQLC generated
│   │       ├── querier.go             # SQLC generated
│   │       ├── users.sql.go           # SQLC generated
│   │       ├── courses.sql.go         # SQLC generated
│   │       └── ...                    # Other SQLC files
│   │
│   ├── service/
│   │   ├── auth/
│   │   │   ├── service.go             # Auth service interface
│   │   │   ├── implementation.go      # Auth implementation
│   │   │   ├── jwt.go                 # JWT helpers
│   │   │   ├── password.go            # Password utilities
│   │   │   └── session.go             # Session management
│   │   │
│   │   ├── user/
│   │   │   ├── service.go             # User service interface
│   │   │   ├── implementation.go      # User implementation
│   │   │   └── profile.go             # Profile management
│   │   │
│   │   ├── course/
│   │   │   ├── service.go             # Course service interface
│   │   │   ├── implementation.go      # Course implementation
│   │   │   ├── enrollment.go          # Enrollment logic
│   │   │   └── progress.go            # Progress tracking
│   │   │
│   │   ├── assessment/
│   │   │   ├── service.go             # Assessment service interface
│   │   │   ├── implementation.go      # Assessment implementation
│   │   │   ├── quiz.go                # Quiz logic
│   │   │   ├── grading.go             # Grading engine
│   │   │   └── randomizer.go          # Question randomization
│   │   │
│   │   ├── forum/
│   │   │   ├── service.go             # Forum service interface
│   │   │   ├── implementation.go      # Forum implementation
│   │   │   └── moderation.go          # Content moderation
│   │   │
│   │   ├── analytics/
│   │   │   ├── service.go             # Analytics service interface
│   │   │   ├── implementation.go      # Analytics implementation
│   │   │   ├── aggregator.go          # Data aggregation
│   │   │   └── reports.go             # Report generation
│   │   │
│   │   ├── notification/
│   │   │   ├── service.go             # Notification service interface
│   │   │   ├── implementation.go      # Notification implementation
│   │   │   ├── email.go               # Email notifications
│   │   │   └── push.go                # Push notifications
│   │   │
│   │   └── storage/
│   │       ├── service.go             # Storage service interface
│   │       ├── implementation.go      # Storage implementation
│   │       ├── s3.go                  # S3 implementation
│   │       └── local.go               # Local storage
│   │
│   ├── handler/
│   │   ├── auth/
│   │   │   ├── handler.go             # Auth handler setup
│   │   │   ├── login.go               # Login endpoint
│   │   │   ├── register.go            # Registration endpoint
│   │   │   ├── refresh.go             # Token refresh
│   │   │   ├── logout.go              # Logout endpoint
│   │   │   └── password_reset.go      # Password reset
│   │   │
│   │   ├── user/
│   │   │   ├── handler.go             # User handler setup
│   │   │   ├── profile.go             # Profile endpoints
│   │   │   ├── settings.go            # User settings
│   │   │   └── admin.go               # Admin user endpoints
│   │   │
│   │   ├── course/
│   │   │   ├── handler.go             # Course handler setup
│   │   │   ├── crud.go                # CRUD operations
│   │   │   ├── enrollment.go          # Enrollment endpoints
│   │   │   ├── module.go              # Module endpoints
│   │   │   └── lesson.go              # Lesson endpoints
│   │   │
│   │   ├── assessment/
│   │   │   ├── handler.go             # Assessment handler setup
│   │   │   ├── quiz.go                # Quiz endpoints
│   │   │   ├── attempt.go             # Quiz attempt endpoints
│   │   │   └── grading.go             # Grading endpoints
│   │   │
│   │   ├── forum/
│   │   │   ├── handler.go             # Forum handler setup
│   │   │   ├── thread.go              # Thread endpoints
│   │   │   ├── post.go                # Post endpoints
│   │   │   └── moderation.go          # Moderation endpoints
│   │   │
│   │   ├── analytics/
│   │   │   ├── handler.go             # Analytics handler setup
│   │   │   ├── course_analytics.go    # Course analytics endpoints
│   │   │   ├── student_analytics.go   # Student analytics endpoints
│   │   │   └── reports.go             # Report endpoints
│   │   │
│   │   └── admin/
│   │       ├── handler.go             # Admin handler setup
│   │       ├── users.go               # User management
│   │       ├── system.go              # System settings
│   │       └── permissions.go         # Permission management
│   │
│   ├── middleware/
│   │   ├── auth.go                    # Authentication middleware
│   │   ├── permission.go              # Permission checking
│   │   ├── ratelimit.go               # Rate limiting
│   │   ├── logging.go                 # Request logging
│   │   ├── cors.go                    # CORS handling
│   │   ├── security.go                # Security headers
│   │   └── recovery.go                # Panic recovery
│   │
│   ├── repository/
│   │   ├── interfaces.go              # Repository interfaces
│   │   ├── user_repo.go               # User repository
│   │   ├── course_repo.go             # Course repository
│   │   ├── assessment_repo.go         # Assessment repository
│   │   └── transaction.go             # Transaction management
│   │
│   ├── model/
│   │   ├── user.go                    # User domain models
│   │   ├── course.go                  # Course domain models
│   │   ├── assessment.go              # Assessment domain models
│   │   ├── forum.go                   # Forum domain models
│   │   └── common.go                  # Common models
│   │
│   ├── dto/
│   │   ├── auth.go                    # Auth DTOs
│   │   ├── user.go                    # User DTOs
│   │   ├── course.go                  # Course DTOs
│   │   ├── assessment.go              # Assessment DTOs
│   │   ├── pagination.go              # Pagination DTOs
│   │   └── response.go                # Response wrappers
│   │
│   ├── permission/
│   │   ├── permission.go              # Permission definitions
│   │   ├── enforcer.go                # Casbin enforcer
│   │   └── policies.go                # Policy management
│   │
│   ├── worker/
│   │   ├── worker.go                  # Worker setup
│   │   ├── tasks/
│   │   │   ├── email.go               # Email tasks
│   │   │   ├── analytics.go           # Analytics aggregation
│   │   │   ├── notifications.go       # Notification tasks
│   │   │   └── cleanup.go             # Cleanup tasks
│   │   └── scheduler.go               # Task scheduler
│   │
│   └── util/
│       ├── validator/
│       │   ├── validator.go           # Custom validators
│       │   └── rules.go               # Validation rules
│       ├── response/
│       │   ├── json.go                # JSON response helpers
│       │   └── error.go               # Error responses
│       ├── logger/
│       │   ├── logger.go              # Logger setup
│       │   └── context.go             # Context logging
│       └── helper/
│           ├── pagination.go          # Pagination helpers
│           ├── slug.go                # Slug generation
│           └── random.go              # Random generators
│
├── db/
│   ├── migrations/
│   │   ├── 00001_init_schema.sql
│   │   ├── 00002_create_users_table.sql
│   │   ├── 00003_create_auth_tables.sql
│   │   ├── 00004_create_permissions_tables.sql
│   │   ├── 00005_create_courses_table.sql
│   │   ├── 00006_create_modules_lessons_tables.sql
│   │   ├── 00007_create_enrollments_table.sql
│   │   ├── 00008_create_assessments_tables.sql
│   │   ├── 00009_create_forum_tables.sql
│   │   ├── 00010_create_analytics_tables.sql
│   │   ├── 00011_create_notifications_table.sql
│   │   ├── 00012_add_indexes.sql
│   │   └── 00013_add_triggers.sql
│   │
│   ├── queries/
│   │   ├── users.sql
│   │   ├── auth.sql
│   │   ├── permissions.sql
│   │   ├── courses.sql
│   │   ├── modules.sql
│   │   ├── lessons.sql
│   │   ├── enrollments.sql
│   │   ├── assessments.sql
│   │   ├── forums.sql
│   │   ├── analytics.sql
│   │   └── notifications.sql
│   │
│   └── seeds/
│       ├── 01_users.sql
│       ├── 02_permissions.sql
│       ├── 03_courses.sql
│       └── 04_test_data.sql
│
├── pkg/
│   ├── cache/
│   │   ├── cache.go                   # Cache interface
│   │   ├── memory.go                  # In-memory implementation
│   │   └── redis.go                   # Redis implementation
│   │
│   ├── email/
│   │   ├── sender.go                  # Email sender interface
│   │   ├── smtp.go                    # SMTP implementation
│   │   └── templates/                 # Email templates
│   │
│   └── security/
│       ├── csrf.go                    # CSRF protection
│       └── sanitizer.go               # Input sanitization
│
├── web/
│   ├── static/                        # Static files
│   └── templates/                     # HTML templates
│
├── scripts/
│   ├── setup.sh                       # Project setup
│   ├── generate.sh                    # Code generation
│   └── docker-init.sh                 # Docker initialization
│
├── test/
│   ├── integration/                   # Integration tests
│   ├── fixtures/                      # Test fixtures
│   └── mocks/                         # Generated mocks
│
├── docs/
│   ├── api/                          # API documentation
│   └── architecture/                  # Architecture docs
│
├── .github/
│   └── workflows/
│       ├── test.yml
│       └── deploy.yml
│
├── build/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── configs/
│   ├── .env.example
│   ├── casbin_model.conf
│   └── casbin_policy.csv
│
├── Makefile
├── sqlc.yaml
├── .golangci.yml
├── go.mod
├── go.sum
└── README.md
```

## SQLC Configuration (sqlc.yaml)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries"
    schema: "db/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/infrastructure/database/sqlc"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_exported_queries: false
        emit_result_struct_pointers: false
        emit_params_struct_pointers: false
        emit_methods_with_db_argument: false
        emit_pointers_for_null_types: true
        emit_enum_valid_method: true
        emit_all_enum_values: true
        json_tags_case_style: "camel"
        output_models_file_name: "models.go"
        output_querier_file_name: "querier.go"
        output_copyfrom_file_name: "copyfrom.go"
        query_parameter_limit: 1000
        omit_unused_structs: true
    rules:
      - sqlc/db-prepare
    overrides:
      - db_type: "uuid"
        go_type: "github.com/google/uuid.UUID"
      - db_type: "timestamptz"
        go_type: "time.Time"
      - db_type: "jsonb"
        go_type: 
          type: "json.RawMessage"
```
## Complete Database Schema
### Core Tables

```sql
-- Users and Authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    email_verification_token VARCHAR(255),
    email_verified_at TIMESTAMP,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200),
    avatar_url VARCHAR(500),
    bio TEXT,
    phone VARCHAR(20),
    date_of_birth DATE,
    gender VARCHAR(20),
    country VARCHAR(2),
    timezone VARCHAR(50) DEFAULT 'UTC',
    preferred_language VARCHAR(10) DEFAULT 'en',
    is_active BOOLEAN DEFAULT true,
    suspended_at TIMESTAMP,
    suspended_reason TEXT,
    last_login_at TIMESTAMP,
    login_count INTEGER DEFAULT 0,
    failed_login_attempts INTEGER DEFAULT 0,
    failed_login_locked_until TIMESTAMP,
    password_changed_at TIMESTAMP,
    must_change_password BOOLEAN DEFAULT false,
    two_factor_enabled BOOLEAN DEFAULT false,
    two_factor_secret VARCHAR(255),
    backup_codes TEXT[], -- Array of encrypted backup codes
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP, -- Soft delete
    CONSTRAINT chk_email CHECK (email ~* '^.+@.+\..+$'),
    INDEX idx_users_email (email),
    INDEX idx_users_active (is_active, deleted_at),
    INDEX idx_users_created (created_at)
);

-- User Sessions
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) UNIQUE NOT NULL,
    access_token_hash VARCHAR(255),
    device_name VARCHAR(255),
    device_type VARCHAR(50), -- mobile, desktop, tablet
    browser VARCHAR(100),
    browser_version VARCHAR(50),
    os VARCHAR(100),
    os_version VARCHAR(50),
    ip_address INET,
    location JSONB, -- {country, city, region}
    is_active BOOLEAN DEFAULT true,
    last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    revoked_reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sessions_user (user_id, is_active),
    INDEX idx_sessions_token (refresh_token_hash),
    INDEX idx_sessions_expires (expires_at)
);

-- Password Reset Tokens
CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_password_resets_token (token_hash),
    INDEX idx_password_resets_user (user_id)
);

-- Groups (Roles)
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT false, -- Cannot be deleted
    priority INTEGER DEFAULT 0, -- For permission precedence
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_groups_name (name)
);

-- Permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL, -- courses, users, analytics
    action VARCHAR(100) NOT NULL,   -- create, read, update, delete
    scope VARCHAR(50) NOT NULL,     -- own, assigned, all
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action, scope),
    INDEX idx_permissions_resource (resource)
);

-- Group Permissions
CREATE TABLE group_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    constraints JSONB DEFAULT '{}', -- Additional constraints
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    granted_by UUID REFERENCES users(id),
    UNIQUE(group_id, permission_id),
    INDEX idx_group_permissions_group (group_id),
    INDEX idx_group_permissions_permission (permission_id)
);

-- User Groups
CREATE TABLE user_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(id),
    expires_at TIMESTAMP, -- For temporary assignments
    UNIQUE(user_id, group_id),
    INDEX idx_user_groups_user (user_id),
    INDEX idx_user_groups_group (group_id)
);

-- Courses
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
    deleted_at TIMESTAMP,
    INDEX idx_courses_instructor (instructor_id),
    INDEX idx_courses_published (is_published, deleted_at),
    INDEX idx_courses_category (category, sub_category),
    INDEX idx_courses_slug (slug),
    INDEX idx_courses_search (title, description)
);

-- Course Co-instructors & Teaching Assistants
CREATE TABLE course_staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- co_instructor, teaching_assistant
    permissions JSONB DEFAULT '{}', -- Specific permissions for this course
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(id),
    UNIQUE(course_id, user_id),
    INDEX idx_course_staff_course (course_id),
    INDEX idx_course_staff_user (user_id)
);

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
    UNIQUE(course_id, order_index),
    INDEX idx_modules_course (course_id, order_index)
);

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
    UNIQUE(module_id, order_index),
    INDEX idx_lessons_module (module_id, order_index)
);

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
    UNIQUE(user_id, course_id),
    INDEX idx_enrollments_user (user_id, status),
    INDEX idx_enrollments_course (course_id, status),
    INDEX idx_enrollments_completed (completed_at)
);

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
    UNIQUE(user_id, course_id),
    INDEX idx_enrollment_requests_course (course_id, status)
);

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_access_codes_code (code),
    INDEX idx_access_codes_course (course_id)
);

-- Quizzes
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    instructions TEXT,
    quiz_type VARCHAR(50) DEFAULT 'graded', -- graded, practice, survey
    order_index INTEGER NOT NULL,
    is_published BOOLEAN DEFAULT true,
    available_from TIMESTAMP,
    available_until TIMESTAMP,
    time_limit_minutes INTEGER,
    attempt_limit INTEGER DEFAULT 1,
    passing_score DECIMAL(5,2) DEFAULT 70.00,
    total_points INTEGER DEFAULT 0,
    randomize_questions BOOLEAN DEFAULT false,
    randomize_answers BOOLEAN DEFAULT false,
    questions_per_page INTEGER DEFAULT 1,
    show_correct_answers VARCHAR(50) DEFAULT 'after_submission', -- never, after_submission, after_deadline
    allow_back_navigation BOOLEAN DEFAULT true,
    required_for_completion BOOLEAN DEFAULT true,
    weight_percentage DECIMAL(5,2), -- Weight in final grade
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(module_id, order_index),
    INDEX idx_quizzes_module (module_id)
);

-- Question Bank (Reusable questions)
CREATE TABLE question_bank (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    tags TEXT[],
    difficulty VARCHAR(50), -- easy, medium, hard
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL, -- multiple_choice, true_false, fill_blank, matching, essay
    explanation TEXT,
    hints TEXT[],
    points INTEGER DEFAULT 1,
    time_estimate_seconds INTEGER,
    usage_count INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_question_bank_creator (created_by),
    INDEX idx_question_bank_category (category),
    INDEX idx_question_bank_type (question_type)
);

-- Quiz Questions
CREATE TABLE quiz_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    question_bank_id UUID REFERENCES question_bank(id), -- If from bank
    order_index INTEGER NOT NULL,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL,
    required BOOLEAN DEFAULT true,
    points INTEGER DEFAULT 1,
    negative_points INTEGER DEFAULT 0, -- For negative marking
    explanation TEXT,
    hints TEXT[],
    time_limit_seconds INTEGER, -- Per-question time limit
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(quiz_id, order_index),
    INDEX idx_quiz_questions_quiz (quiz_id, order_index)
);

-- Answer Options (for multiple choice, matching, etc.)
CREATE TABLE answer_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
    option_text TEXT NOT NULL,
    option_value TEXT, -- For matching questions
    is_correct BOOLEAN DEFAULT false,
    explanation TEXT, -- Why this option is correct/incorrect
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(question_id, order_index),
    INDEX idx_answer_options_question (question_id)
);

-- Quiz Attempts
CREATE TABLE quiz_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    attempt_number INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(50) DEFAULT 'in_progress', -- in_progress, submitted, graded, abandoned
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMP,
    graded_at TIMESTAMP,
    time_spent_seconds INTEGER DEFAULT 0,
    score DECIMAL(5,2),
    points_earned INTEGER DEFAULT 0,
    passed BOOLEAN,
    ip_address INET,
    browser_info JSONB,
    flagged_for_review BOOLEAN DEFAULT false, -- Cheating detection
    review_notes TEXT,
    graded_by UUID REFERENCES users(id),
    INDEX idx_quiz_attempts_user (user_id),
    INDEX idx_quiz_attempts_quiz (quiz_id),
    INDEX idx_quiz_attempts_status (status)
);

-- Student Answers
CREATE TABLE student_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES quiz_questions(id),
    answer_text TEXT, -- For text/essay answers
    selected_options UUID[], -- Array of answer_option IDs
    is_correct BOOLEAN,
    points_earned DECIMAL(5,2) DEFAULT 0,
    time_spent_seconds INTEGER DEFAULT 0,
    marked_for_review BOOLEAN DEFAULT false,
    feedback TEXT, -- Instructor feedback
    graded_at TIMESTAMP,
    graded_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(attempt_id, question_id),
    INDEX idx_student_answers_attempt (attempt_id)
);

-- Lesson Progress
CREATE TABLE lesson_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    enrollment_id UUID NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'not_started', -- not_started, in_progress, completed
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    last_position INTEGER DEFAULT 0, -- For videos (seconds) or text (character position)
    time_spent_seconds INTEGER DEFAULT 0,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    notes TEXT, -- Student notes
    bookmarks JSONB DEFAULT '[]', -- Array of {position, note, created_at}
    UNIQUE(user_id, lesson_id),
    INDEX idx_lesson_progress_user (user_id, status),
    INDEX idx_lesson_progress_lesson (lesson_id),
    INDEX idx_lesson_progress_enrollment (enrollment_id)
);

-- Module Progress (Calculated/Cached)
CREATE TABLE module_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    enrollment_id UUID NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    lessons_completed INTEGER DEFAULT 0,
    lessons_total INTEGER DEFAULT 0,
    quizzes_completed INTEGER DEFAULT 0,
    quizzes_total INTEGER DEFAULT 0,
    average_quiz_score DECIMAL(5,2),
    time_spent_seconds INTEGER DEFAULT 0,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    unlocked_at TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    UNIQUE(user_id, module_id),
    INDEX idx_module_progress_user (user_id),
    INDEX idx_module_progress_enrollment (enrollment_id)
);

-- Forum Threads
CREATE TABLE forum_threads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE, -- Optional module-specific
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE, -- Optional lesson-specific
    author_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    thread_type VARCHAR(50) DEFAULT 'discussion', -- discussion, question, announcement
    is_pinned BOOLEAN DEFAULT false,
    is_locked BOOLEAN DEFAULT false,
    is_anonymous BOOLEAN DEFAULT false,
    tags TEXT[],
    view_count INTEGER DEFAULT 0,
    reply_count INTEGER DEFAULT 0,
    last_reply_at TIMESTAMP,
    last_reply_by UUID REFERENCES users(id),
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    is_answered BOOLEAN DEFAULT false, -- For questions
    best_answer_id UUID, -- References forum_posts
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id),
    INDEX idx_forum_threads_course (course_id, deleted_at),
    INDEX idx_forum_threads_author (author_id),
    INDEX idx_forum_threads_type (thread_type),
    FULLTEXT INDEX idx_forum_threads_search (title, content)
);

-- Forum Posts
CREATE TABLE forum_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id UUID NOT NULL REFERENCES forum_threads(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    is_anonymous BOOLEAN DEFAULT false,
    is_instructor_response BOOLEAN DEFAULT false,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    is_solution BOOLEAN DEFAULT false,
    edited_at TIMESTAMP,
    edited_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id),
    deleted_reason TEXT,
    INDEX idx_forum_posts_thread (thread_id, deleted_at),
    INDEX idx_forum_posts_author (author_id),
    INDEX idx_forum_posts_parent (parent_id)
);

-- Forum Votes
CREATE TABLE forum_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    votable_type VARCHAR(50) NOT NULL, -- thread, post
    votable_id UUID NOT NULL,
    vote_type VARCHAR(10) NOT NULL, -- up, down
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, votable_type, votable_id),
    INDEX idx_forum_votes_user (user_id),
    INDEX idx_forum_votes_votable (votable_type, votable_id)
);

-- Forum Reports
CREATE TABLE forum_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id UUID NOT NULL REFERENCES users(id),
    content_type VARCHAR(50) NOT NULL, -- thread, post
    content_id UUID NOT NULL,
    reason VARCHAR(100) NOT NULL, -- spam, inappropriate, harassment, etc.
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- pending, reviewed, resolved, dismissed
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP,
    action_taken VARCHAR(100), -- removed, warned, no_action
    moderator_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_forum_reports_status (status),
    INDEX idx_forum_reports_content (content_type, content_id)
);

-- Notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL, -- enrollment_accepted, quiz_graded, new_reply, etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB DEFAULT '{}', -- Additional data
    priority VARCHAR(20) DEFAULT 'normal', -- low, normal, high, urgent
    read_at TIMESTAMP,
    clicked_at TIMESTAMP,
    action_url VARCHAR(500),
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_notifications_user (user_id, read_at),
    INDEX idx_notifications_created (created_at),
    INDEX idx_notifications_type (type)
);

-- Notification Preferences
CREATE TABLE notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email_enabled BOOLEAN DEFAULT true,
    email_frequency VARCHAR(50) DEFAULT 'immediate', -- immediate, daily, weekly
    push_enabled BOOLEAN DEFAULT true,
    sms_enabled BOOLEAN DEFAULT false,
    preferences JSONB DEFAULT '{}', -- {quiz_graded: true, new_announcement: true, ...}
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_notification_preferences_user (user_id)
);

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_analytics_events_user (user_id, created_at),
    INDEX idx_analytics_events_course (course_id, created_at),
    INDEX idx_analytics_events_type (event_type, created_at),
    INDEX idx_analytics_events_created (created_at)
);

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
    UNIQUE(course_id, date),
    INDEX idx_course_analytics_course (course_id, date)
);

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_system_logs_user (user_id, created_at),
    INDEX idx_system_logs_action (action, created_at),
    INDEX idx_system_logs_created (created_at)
);

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_audit_logs_user (user_id, created_at),
    INDEX idx_audit_logs_resource (resource_type, resource_id),
    INDEX idx_audit_logs_created (created_at)
);

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
    deleted_at TIMESTAMP,
    INDEX idx_file_uploads_user (uploaded_by),
    INDEX idx_file_uploads_type (file_type)
);

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
    UNIQUE(course_id, user_id),
    INDEX idx_course_ratings_course (course_id, rating),
    INDEX idx_course_ratings_user (user_id)
);

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
    revoked_reason TEXT,
    INDEX idx_certificates_enrollment (enrollment_id),
    INDEX idx_certificates_number (certificate_number)
);

-- Announcements
CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal', -- low, normal, high, urgent
    is_published BOOLEAN DEFAULT true,
    publish_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    target_audience JSONB DEFAULT '{}', -- {groups: [], users: []}
    read_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_announcements_course (course_id, is_published),
    INDEX idx_announcements_publish (publish_at)
);

-- User Announcement Reads
CREATE TABLE announcement_reads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    announcement_id UUID NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(announcement_id, user_id),
    INDEX idx_announcement_reads_user (user_id)
);

-- Backup and Archive tables
CREATE TABLE data_archives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    archive_type VARCHAR(50) NOT NULL, -- user_data, course_data, full_backup
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    file_path VARCHAR(500),
    file_size BIGINT,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    expires_at TIMESTAMP,
    download_count INTEGER DEFAULT 0,
    error_message TEXT,
    INDEX idx_data_archives_user (user_id, status)
);
```

### Database Functions and Triggers

```sql
-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Repeat for other tables...

-- Function to calculate course rating
CREATE OR REPLACE FUNCTION update_course_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE courses
    SET 
        rating_average = (SELECT AVG(rating) FROM course_ratings WHERE course_id = NEW.course_id),
        rating_count = (SELECT COUNT(*) FROM course_ratings WHERE course_id = NEW.course_id)
    WHERE id = NEW.course_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_course_rating_trigger
AFTER INSERT OR UPDATE OR DELETE ON course_ratings
FOR EACH ROW EXECUTE FUNCTION update_course_rating();

-- Function to update enrollment count
CREATE OR REPLACE FUNCTION update_enrollment_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE courses SET enrolled_count = enrolled_count + 1 WHERE id = NEW.course_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE courses SET enrolled_count = enrolled_count - 1 WHERE id = OLD.course_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_enrollment_count_trigger
AFTER INSERT OR DELETE ON enrollments
FOR EACH ROW EXECUTE FUNCTION update_enrollment_count();
```

## Authentication & Session Tasks

### POST /api/v1/auth/register
**Task: User Registration**
1. **Input Validation**
   - Validate email format and uniqueness
   - Check password strength (min 8 chars, uppercase, lowercase, number, special char)
   - Validate required fields (firstName, lastName, email, password)
   - Check if email domain is allowed (optional whitelist)

2. **User Creation**
   - Hash password using bcrypt (cost factor 12)
   - Generate email verification token (UUID)
   - Create user record with `email_verified = false`
   - Assign default 'student' group

3. **Post-Registration**
   - Send verification email with token link
   - Create welcome notification
   - Log registration event
   - Return success message (no auto-login)

**Error Cases:**
- Email already exists
- Weak password
- Invalid email format
- Database errors

### POST /api/v1/auth/login
**Task: User Login**
1. **Pre-Authentication Checks**
   - Check if user exists by email
   - Verify account is not locked (failed attempts)
   - Check if account is active (not suspended)
   - Verify email is confirmed

2. **Authentication**
   - Compare password with bcrypt hash
   - If failed, increment failed_login_attempts
   - Lock account after 5 failed attempts (30 min)
   - Reset failed attempts on success

3. **Session Creation**
   - Generate access token (15 min) with claims: {userId, permissions, sessionId}
   - Generate refresh token (7 days)
   - Store session with device info, IP, location
   - Update last_login_at, increment login_count

4. **Response**
   - Return tokens, user info, permissions
   - Set secure, httpOnly cookies for tokens
   - Include CSRF token if using cookies

**Error Cases:**
- Invalid credentials
- Account locked
- Email not verified
- Account suspended

### POST /api/v1/auth/logout
**Task: User Logout**
1. **Token Validation**
   - Validate access token from header/cookie
   - Extract session ID from token

2. **Session Termination**
   - Mark session as revoked
   - Set revoked_reason = "user_logout"
   - Clear any cached permissions

3. **Cleanup**
   - Clear cookies if used
   - Log logout event
   - Return success

### POST /api/v1/auth/refresh
**Task: Refresh Access Token**
1. **Refresh Token Validation**
   - Validate refresh token format
   - Check if token exists and not revoked
   - Verify token hasn't expired
   - Check associated user is still active

2. **Token Generation**
   - Generate new access token
   - Optionally rotate refresh token
   - Update session last_accessed_at

3. **Response**
   - Return new access token
   - Return new refresh token if rotated

**Error Cases:**
- Invalid/expired refresh token
- Session revoked
- User suspended

### POST /api/v1/auth/forgot-password
**Task: Password Reset Request**
1. **User Lookup**
   - Find user by email
   - Don't reveal if email exists (security)

2. **Token Generation**
   - Generate secure reset token
   - Hash token before storing
   - Set expiry (1 hour)
   - Invalidate any existing tokens

3. **Email Dispatch**
   - Send reset email with link
   - Include user's name
   - Add security notice

4. **Response**
   - Always return success (security)
   - Log password reset request

### POST /api/v1/auth/reset-password
**Task: Password Reset Completion**
1. **Token Validation**
   - Validate token format
   - Check if token exists and not used
   - Verify not expired
   - Find associated user

2. **Password Update**
   - Validate new password strength
   - Hash new password
   - Update user password
   - Mark token as used

3. **Post-Reset**
   - Revoke all user sessions (security)
   - Send confirmation email
   - Log password change
   - Force re-login

### GET /api/v1/auth/sessions
**Task: List User Sessions**
1. **Authorization**
   - Validate user is authenticated
   - Can only view own sessions

2. **Data Retrieval**
   - Get all active sessions for user
   - Include device info, location, last access

3. **Response**
   - Return session list
   - Highlight current session

### DELETE /api/v1/auth/sessions/:sessionId
**Task: Revoke Specific Session**
1. **Authorization**
   - Validate user owns the session
   - Cannot revoke current session

2. **Revocation**
   - Mark session as revoked
   - Set revoked_reason = "user_revoked"
   - Log revocation event

## User Management Tasks

### GET /api/v1/users/me
**Task: Get Current User Profile**
1. **Authentication**
   - Validate access token
   - Extract user ID from token

2. **Data Retrieval**
   - Get user details
   - Include groups and permissions
   - Get notification preferences
   - Calculate storage usage

3. **Response Formatting**
   - Exclude sensitive fields (password_hash)
   - Include computed fields (full_name)
   - Add avatar URL with CDN

### PUT /api/v1/users/me
**Task: Update User Profile**
1. **Input Validation**
   - Validate allowed fields only
   - Check email uniqueness if changed
   - Validate phone format if provided

2. **Update Process**
   - Update allowed fields
   - If email changed, require reverification
   - Log profile changes for audit

3. **Side Effects**
   - Update search indices
   - Clear user cache
   - Send confirmation email for critical changes

### POST /api/v1/users/me/avatar
**Task: Upload Avatar**
1. **File Validation**
   - Check file size (max 5MB)
   - Validate image format (jpg, png, webp)
   - Scan for malware

2. **Image Processing**
   - Resize to standard sizes (50x50, 200x200)
   - Generate webp version
   - Create thumbnail

3. **Storage**
   - Upload to storage service
   - Update user avatar_url
   - Delete old avatar if exists

### GET /api/v1/users (Admin)
**Task: List All Users**
1. **Authorization**
   - Require admin permission
   - Log admin access

2. **Query Building**
   - Apply filters (role, status, created date)
   - Search by name/email
   - Pagination (default 50)
   - Sorting options

3. **Data Enrichment**
   - Include group memberships
   - Add last login info
   - Calculate activity metrics

## Course Management Tasks

### GET /api/v1/courses
**Task: List Courses**
1. **Query Parameters**
   - Filters: category, level, language, price
   - Search: title, description, tags
   - Sorting: popular, newest, rating, price
   - Pagination: page, limit

2. **Authorization Logic**
   - Public: Only published, non-deleted
   - Student: Add enrolled filter option
   - Instructor: Include own drafts
   - Admin: See all including deleted

3. **Data Enrichment**
   - Include instructor info
   - Add enrollment count
   - Calculate duration
   - Include first module preview

### POST /api/v1/courses
**Task: Create Course**
1. **Authorization**
   - Check user has instructor role
   - Verify not exceeded course limit

2. **Validation**
   - Validate required fields
   - Generate unique slug from title
   - Ensure code is unique
   - Validate category exists

3. **Course Setup**
   - Create course in draft status
   - Create default module "Introduction"
   - Set instructor as owner
   - Initialize analytics record

4. **Post-Creation**
   - Send creation notification
   - Log course creation
   - Return with edit URL

### PUT /api/v1/courses/:courseId
**Task: Update Course**
1. **Authorization**
   - Check user owns course or is admin
   - Verify course not deleted

2. **Validation**
   - Cannot change certain fields if published
   - Validate slug uniqueness if changed
   - Check enrollment type change impacts

3. **Update Logic**
   - Update allowed fields
   - Regenerate duration if content changed
   - Update search indices
   - Clear course cache

### POST /api/v1/courses/:courseId/publish
**Task: Publish Course**
1. **Pre-Publication Checks**
   - Verify has at least one module
   - Each module has content
   - All required metadata present
   - Check instructor verified

2. **Publication Process**
   - Set is_published = true
   - Set published_at timestamp
   - Generate course thumbnail
   - Index for search

3. **Notifications**
   - Notify enrolled students (if any)
   - Update instructor dashboard
   - Send to course catalog

## Module & Lesson Tasks

### POST /api/v1/courses/:courseId/modules
**Task: Create Module**
1. **Authorization**
   - Verify course ownership
   - Check course not archived

2. **Validation**
   - Validate title uniqueness in course
   - Check order_index not duplicate

3. **Creation**
   - Insert module with auto-increment order
   - Update course duration estimate
   - Log module creation

### POST /api/v1/modules/:moduleId/lessons
**Task: Create Lesson**
1. **Authorization**
   - Verify module ownership through course

2. **Content Handling**
   - Video: Process upload, generate thumbnail
   - Text: Sanitize HTML, extract reading time
   - PDF: Upload, extract page count
   - Interactive: Validate embed code

3. **Creation Process**
   - Auto-assign order_index
   - Calculate duration based on content
   - Update module duration
   - Create progress tracking records

### POST /api/v1/lessons/:lessonId/start
**Task: Start Lesson**
1. **Authorization**
   - Verify enrollment in course
   - Check module unlocked

2. **Progress Tracking**
   - Create/update progress record
   - Set started_at if first time
   - Start time tracking

3. **Analytics**
   - Log lesson_started event
   - Track device/browser
   - Update daily active users

### POST /api/v1/lessons/:lessonId/complete
**Task: Complete Lesson**
1. **Validation**
   - Check if started
   - Verify minimum time spent
   - Validate completion criteria

2. **Progress Update**
   - Mark lesson completed
   - Update module progress
   - Check if module completed
   - Update course progress

3. **Unlocking Logic**
   - Check next module prerequisites
   - Unlock if conditions met
   - Send notification if course completed

## Assessment System Tasks

### POST /api/v1/modules/:moduleId/quizzes
**Task: Create Quiz**
1. **Authorization**
   - Verify module ownership

2. **Quiz Setup**
   - Validate quiz settings
   - Set default time limits
   - Initialize question array

3. **Configuration**
   - Set grading scheme
   - Configure attempt rules
   - Set availability window

### POST /api/v1/quizzes/:quizId/questions
**Task: Add Question**
1. **Question Validation**
   - Validate question type
   - Check required fields per type
   - Validate point value

2. **Answer Setup**
   - Multiple choice: Validate options, mark correct
   - True/False: Set correct answer
   - Essay: Set grading rubric
   - Fill blank: Set acceptable answers

3. **Quiz Update**
   - Update total points
   - Reorder if needed
   - Update duration estimate

### POST /api/v1/quizzes/:quizId/start
**Task: Start Quiz Attempt**
1. **Eligibility Check**
   - Verify enrollment
   - Check attempt limit
   - Verify availability window
   - Check prerequisites

2. **Attempt Creation**
   - Create attempt record
   - Generate question order (if random)
   - Randomize answer options
   - Start timer

3. **Response**
   - Return questions without answers
   - Include time limit
   - Set attempt token

### POST /api/v1/attempts/:attemptId/submit
**Task: Submit Quiz**
1. **Validation**
   - Verify attempt ownership
   - Check not already submitted
   - Validate time limit

2. **Answer Processing**
   - Save all provided answers
   - Auto-grade objective questions
   - Calculate preliminary score
   - Flag for manual grading if needed

3. **Completion**
   - Update attempt status
   - Calculate final score
   - Update progress
   - Send notification
   - Generate certificate if passed

### POST /api/v1/attempts/:attemptId/grade
**Task: Manual Grading**
1. **Authorization**
   - Verify grader permission
   - Check attempt needs grading

2. **Grading Process**
   - Grade each essay/open question
   - Provide feedback
   - Calculate total score

3. **Finalization**
   - Update attempt score
   - Mark as graded
   - Notify student
   - Update gradebook

## Enrollment Tasks

### POST /api/v1/enrollments
**Task: Enroll in Course**
1. **Eligibility**
   - Check course published
   - Verify enrollment open
   - Check capacity
   - Validate prerequisites

2. **Enrollment Type**
   - Open: Direct enrollment
   - Approval: Create request
   - Access Code: Validate code
   - Payment: Verify payment

3. **Enrollment Process**
   - Create enrollment record
   - Initialize progress tracking
   - Grant course permissions
   - Send welcome email

### GET /api/v1/courses/:courseId/enrollments
**Task: List Course Students**
1. **Authorization**
   - Verify instructor/TA permission

2. **Data Retrieval**
   - Get enrolled students
   - Include progress data
   - Add last activity
   - Filter by status

3. **Export Options**
   - CSV export
   - Include grades
   - Contact information

## Forum System Tasks

### POST /api/v1/courses/:courseId/threads
**Task: Create Thread**
1. **Authorization**
   - Verify enrollment
   - Check forum permissions

2. **Content Processing**
   - Sanitize HTML content
   - Extract mentions
   - Process tags

3. **Thread Creation**
   - Create thread record
   - Send notifications to mentioned
   - Index for search
   - Award participation points

### POST /api/v1/threads/:threadId/posts
**Task: Reply to Thread**
1. **Authorization**
   - Check thread not locked
   - Verify enrollment

2. **Reply Processing**
   - Link to parent if nested
   - Process mentions
   - Check for spam

3. **Notifications**
   - Notify thread author
   - Notify mentioned users
   - Update thread activity

### POST /api/v1/posts/:postId/report
**Task: Report Content**
1. **Report Creation**
   - Validate reason
   - Create report record
   - Capture context

2. **Moderation Queue**
   - Add to mod queue
   - Notify moderators
   - Set priority by severity

3. **Auto-Actions**
   - Hide if threshold reached
   - Flag user if repeat offender

## Analytics Tasks

### GET /api/v1/analytics/courses/:courseId/overview
**Task: Course Analytics Dashboard**
1. **Authorization**
   - Verify instructor permission

2. **Metrics Calculation**
   - Enrollment trends
   - Completion rates
   - Average progress
   - Time spent distribution

3. **Data Aggregation**
   - Daily/weekly/monthly views
   - Module performance
   - Quiz statistics
   - Forum engagement

### POST /api/v1/reports/generate
**Task: Generate Custom Report**
1. **Report Configuration**
   - Select metrics
   - Set date range
   - Choose grouping

2. **Data Processing**
   - Run aggregation queries
   - Calculate trends
   - Generate visualizations

3. **Report Delivery**
   - Generate PDF/Excel
   - Email when ready
   - Store for download

## Administration Tasks

### PUT /api/v1/admin/groups/:groupId/permissions
**Task: Update Group Permissions**
1. **Authorization**
   - Require super admin
   - Log permission changes

2. **Validation**
   - Validate permission IDs
   - Check for conflicts
   - Prevent lockout scenarios

3. **Update Process**
   - Update permissions
   - Clear permission cache
   - Notify affected users
   - Create audit log

### POST /api/v1/admin/backup
**Task: System Backup**
1. **Backup Initiation**
   - Verify admin permission
   - Check storage space

2. **Backup Process**
   - Create database dump
   - Archive uploaded files
   - Export configurations

3. **Completion**
   - Encrypt backup
   - Upload to secure storage
   - Log backup details
   - Send confirmation