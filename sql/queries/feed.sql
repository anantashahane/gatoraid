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
