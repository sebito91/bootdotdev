// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, url, description, published_at, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	FeedID      uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.FeedID,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsByUser = `-- name: GetPostsByUser :many
select id, created_at, updated_at, title, url, description, published_at, feed_id from posts where feed_id = $1
`

func (q *Queries) GetPostsByUser(ctx context.Context, feedID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUser, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}