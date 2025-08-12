-- name: CreateSession :one
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;