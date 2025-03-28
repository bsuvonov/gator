// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getFeedIdByUrl = `-- name: GetFeedIdByUrl :one
SELECT id FROM feeds where url = $1
`

func (q *Queries) GetFeedIdByUrl(ctx context.Context, url string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getFeedIdByUrl, url)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT feeds.name, url, users.name from feeds inner join users on feeds.user_id = users.id
`

type GetFeedsRow struct {
	Name   string
	Url    string
	Name_2 string
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(&i.Name, &i.Url, &i.Name_2); err != nil {
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

const insertFeed = `-- name: InsertFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type InsertFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) InsertFeed(ctx context.Context, arg InsertFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, insertFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :one
UPDATE feeds SET last_fetched_at = $2, updated_at = $2 WHERE id = $1 RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type MarkFeedFetchedParams struct {
	ID            uuid.UUID
	LastFetchedAt sql.NullTime
}

func (q *Queries) MarkFeedFetched(ctx context.Context, arg MarkFeedFetchedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, markFeedFetched, arg.ID, arg.LastFetchedAt)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}
