-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING 
    feed_follows.id,
    feed_follows.user_id,
    feed_follows.feed_id,
    feed_follows.created_at,
    feed_follows.updated_at,
    (SELECT name FROM users WHERE id = feed_follows.user_id) AS user_name,
    (SELECT name FROM feeds WHERE id = feed_follows.feed_id) AS feed_name;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = (
    SELECT id FROM feeds WHERE url = $2
);

-- name: GetFeedFollowsForUser :many
SELECT 
    ff.id,
    ff.user_id,
    ff.feed_id,
    ff.created_at,
    ff.updated_at,
    u.name AS user_name,
    f.name AS feed_name,
    f.url AS feed_url
FROM feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;