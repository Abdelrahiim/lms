-- +goose Up
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

-- +goose Down
-- Drop triggers first
DROP TRIGGER IF EXISTS update_enrollment_count_trigger ON enrollments;
DROP TRIGGER IF EXISTS update_course_rating_trigger ON course_ratings;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop functions
DROP FUNCTION IF EXISTS update_enrollment_count();
DROP FUNCTION IF EXISTS update_course_rating();
DROP FUNCTION IF EXISTS update_updated_at_column();