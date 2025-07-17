-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetMyFeeds :many
SELECT * FROM feeds
WHERE user_id = $1
ORDER BY created_at;

-- name: GetAllFeeds :many
SELECT feeds.name, feeds.url, users.name FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.created_at;

-- name: GetFeed :one
SELECT feeds.name, feeds.url, feeds.id FROM feeds
WHERE feeds.url = $1;

-- name: MarkFeedtoFetch :one
UPDATE feeds
SET updated_at = $1
WHERE feeds.id = $2
RETURNING *;

-- name: GetNextFeedtoFetch :one
SELECT * FROM feeds
ORDER BY updated_at ASC NULLS FIRST
LIMIT 1;
