-- +goose Up
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
    deleted_by UUID REFERENCES users(id)
);

CREATE INDEX idx_forum_threads_course ON forum_threads(course_id, deleted_at);
CREATE INDEX idx_forum_threads_author ON forum_threads(author_id);
CREATE INDEX idx_forum_threads_type ON forum_threads(thread_type);
CREATE INDEX idx_forum_threads_search ON forum_threads USING gin(to_tsvector('english', title || ' ' || content));

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
    deleted_reason TEXT
);

CREATE INDEX idx_forum_posts_thread ON forum_posts(thread_id, deleted_at);
CREATE INDEX idx_forum_posts_author ON forum_posts(author_id);
CREATE INDEX idx_forum_posts_parent ON forum_posts(parent_id);

-- Forum Votes
CREATE TABLE forum_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    votable_type VARCHAR(50) NOT NULL, -- thread, post
    votable_id UUID NOT NULL,
    vote_type VARCHAR(10) NOT NULL, -- up, down
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, votable_type, votable_id)
);

CREATE INDEX idx_forum_votes_user ON forum_votes(user_id);
CREATE INDEX idx_forum_votes_votable ON forum_votes(votable_type, votable_id);

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forum_reports_status ON forum_reports(status);
CREATE INDEX idx_forum_reports_content ON forum_reports(content_type, content_id);

-- +goose Down
DROP TABLE IF EXISTS forum_reports;
DROP TABLE IF EXISTS forum_votes;
DROP TABLE IF EXISTS forum_posts;
DROP TABLE IF EXISTS forum_threads;
