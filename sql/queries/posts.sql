-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
WITH userid AS (
    SELECT id FROM users where users.name = $1
),
feed_ids AS (
    SELECT feed_id FROM feed_follows WHERE user_id = (SELECT id FROM userid)
)
SELECT title, url, description, published_at FROM posts WHERE feed_id IN (SELECT feed_id FROM feed_ids) ORDER BY updated_at ASC LIMIT $2;
