-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetFeedByID :one
select * from feeds where id = $1;

-- name: GetNextFeedsToFetch :many
select * from feeds where last_fetched_at > $1 order by (last_fetched_at, id) limit $2;
