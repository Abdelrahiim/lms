-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash, first_name, last_name, display_name, avatar_url, bio, phone, date_of_birth, gender, country, timezone, preferred_language, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

