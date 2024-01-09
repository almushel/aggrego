-- name: FollowFeed :one
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows WHERE user_id=$1 and id = $2;

-- name: GetUserFollows :many
SELECT * FROM feed_follows
WHERE user_id = $1;