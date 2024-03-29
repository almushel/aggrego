// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: feed_follows.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const followFeed = `-- name: FollowFeed :one
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, user_id, feed_id, created_at, updated_at
`

type FollowFeedParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) FollowFeed(ctx context.Context, arg FollowFeedParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, followFeed, arg.ID, arg.UserID, arg.FeedID)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FeedID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserFollows = `-- name: GetUserFollows :many
SELECT id, user_id, feed_id, created_at, updated_at FROM feed_follows
WHERE user_id = $1
OFFSET $2
LIMIT $3
`

type GetUserFollowsParams struct {
	UserID uuid.UUID
	Offset int32
	Limit  int32
}

func (q *Queries) GetUserFollows(ctx context.Context, arg GetUserFollowsParams) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFollows, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FeedID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const unfollowFeed = `-- name: UnfollowFeed :exec
DELETE FROM feed_follows WHERE user_id=$1 AND id = $2
`

type UnfollowFeedParams struct {
	UserID uuid.UUID
	ID     uuid.UUID
}

func (q *Queries) UnfollowFeed(ctx context.Context, arg UnfollowFeedParams) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, arg.UserID, arg.ID)
	return err
}
