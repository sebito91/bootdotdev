-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: GetFeedFollowsByUser :many
select * from feed_follows where user_id = $1; 
