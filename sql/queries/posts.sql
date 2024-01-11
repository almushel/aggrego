-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostsByUser :many
SELECT *
FROM posts
WHERE feed_id IN (
	SELECT feed_id 
	FROM feed_follows
	WHERE user_id=$1
)
LIMIT $2;