-- name: InsertFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, url, users.name from feeds inner join users on feeds.user_id = users.id;

-- name: GetFeedIdByUrl :one
SELECT id FROM feeds where url = $1;
