-- +goose Up
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
    UNIQUE(module_id, order_index)
);

CREATE INDEX idx_quizzes_module ON quizzes (module_id);

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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_question_bank_creator ON question_bank (created_by);
CREATE INDEX idx_question_bank_category ON question_bank (category);
CREATE INDEX idx_question_bank_type ON question_bank (question_type);

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
    UNIQUE(quiz_id, order_index)
);

CREATE INDEX idx_quiz_questions_quiz ON quiz_questions (quiz_id, order_index);

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
    UNIQUE(question_id, order_index)
);

CREATE INDEX idx_answer_options_question ON answer_options (question_id);

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
    graded_by UUID REFERENCES users(id)
);

CREATE INDEX idx_quiz_attempts_user ON quiz_attempts (user_id);
CREATE INDEX idx_quiz_attempts_quiz ON quiz_attempts (quiz_id);
CREATE INDEX idx_quiz_attempts_status ON quiz_attempts (status);

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
    UNIQUE(attempt_id, question_id)
);

CREATE INDEX idx_student_answers_attempt ON student_answers (attempt_id);

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
    UNIQUE(user_id, lesson_id)
);

CREATE INDEX idx_lesson_progress_user ON lesson_progress (user_id, status);
CREATE INDEX idx_lesson_progress_lesson ON lesson_progress (lesson_id);
CREATE INDEX idx_lesson_progress_enrollment ON lesson_progress (enrollment_id);

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
    UNIQUE(user_id, module_id)
);

CREATE INDEX idx_module_progress_user ON module_progress (user_id);
CREATE INDEX idx_module_progress_enrollment ON module_progress (enrollment_id);

-- +goose Down
DROP TABLE IF EXISTS module_progress;
DROP TABLE IF EXISTS lesson_progress;
DROP TABLE IF EXISTS student_answers;
DROP TABLE IF EXISTS quiz_attempts;
DROP TABLE IF EXISTS answer_options;
DROP TABLE IF EXISTS quiz_questions;
DROP TABLE IF EXISTS question_bank;
DROP TABLE IF EXISTS quizzes;