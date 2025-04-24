-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds ORDER BY created_at DESC;

-- name: GetFeedsWithUsers :many
SELECT 
    f.id,
    f.name,
    f.url,
    f.user_id,
    f.created_at,
    f.updated_at,
    u.name as user_name
FROM feeds f
JOIN users u ON f.user_id = u.id
ORDER BY f.created_at DESC;