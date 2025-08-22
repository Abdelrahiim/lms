-- name: CreateSession :exec
INSERT INTO user_sessions (
        id,
        user_id,
        refresh_token_hash,
        access_token_hash,
        device_name,
        device_type,
        browser,
        browser_version,
        os,
        os_version,
        ip_address,
        location,
        is_active,
        last_accessed_at,
        expires_at,
        revoked_at,
        revoked_reason,
        created_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14,
        $15,
        $16,
        $17,
        $18
    );

-- name: GetSessionByUserID :one
SELECT *
FROM user_sessions
WHERE user_id = $1
    AND ip_address = $2
    AND is_active = TRUE;

-- name: RevokeSession :exec
UPDATE user_sessions
SET is_active = FALSE,
    revoked_at = $1,
    revoked_reason = $2
WHERE id = $3;

-- name: GetActiveSessions :many
SELECT *
FROM user_sessions
WHERE user_id = $1
    AND is_active = $2;

-- name: UpdateSessionLastAccessedAt :exec
UPDATE user_sessions
SET last_accessed_at = $1
WHERE id = $2;

-- name: GetSessionByRefreshToken :one
SELECT *
FROM user_sessions
WHERE refresh_token_hash = $1
    AND is_active = TRUE;
