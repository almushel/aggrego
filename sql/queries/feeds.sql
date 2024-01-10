-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2, $3, $4)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at=CURRENT_TIMESTAMP, last_fetched_at=CURRENT_TIMESTAMP
WHERE id=$1;

-- name: GetStaleFeeds :many
SELECT * FROM feeds
ORDER BY 
CASE 
	WHEN last_fetched_at IS NULL THEN 1
	ELSE 2
END DESC,
last_fetched_at ASC
LIMIT $1;