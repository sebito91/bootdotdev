-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
WHERE last_fetched_at < $1
ORDER BY (last_fetched_at, id)
LIMIT $2;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1
RETURNING *;
