-- name: CreatePost :one
INSERT INTO posts (id, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostsForUser :many
SELECT 
    p.id,
    p.created_at,
    p.updated_at,
    p.title,
    p.url,
    p.description,
    p.published_at,
    p.feed_id,
    f.name AS feed_name
FROM posts p
JOIN feeds f ON p.feed_id = f.id
JOIN feed_follows ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;