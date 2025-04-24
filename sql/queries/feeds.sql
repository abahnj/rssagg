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

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1 LIMIT 1;

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

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST, updated_at
LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;