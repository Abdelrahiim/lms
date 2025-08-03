-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd