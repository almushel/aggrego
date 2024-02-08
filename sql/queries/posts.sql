-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostsByUser :many
SELECT 
	posts.*,
	CAST (
		CASE 
			WHEN 
			EXISTS (select * from post_likes WHERE post_likes.post_id=posts.id)
			THEN 1
			ELSE 0
		END as BOOLEAN
	) AS liked
FROM posts
LEFT JOIN post_likes ON post_likes.user_id=$1
WHERE feed_id IN (
	SELECT feed_id 
	FROM feed_follows
	WHERE user_id=$1
)
OFFSET $2
LIMIT $3;

-- name: GetPostCount :one
SELECT COUNT(*)
FROM posts
WHERE feed_id IN (
	SELECT feed_id 
	FROM feed_follows
	WHERE user_id=$1
);

-- name: LikePost :one
INSERT INTO post_likes (id, user_id, post_id, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UnlikePost :one
DELETE FROM post_likes
WHERE id=$1
RETURNING *;

-- name: GetPostLike :one
SELECT *
FROM post_likes
WHERE id=$1;

-- name: GetLikedPostsByUser :many
SELECT * 
FROM posts
WHERE id IN (
	SELECT post_id
	FROM post_likes
	WHERE user_id=$1
); 